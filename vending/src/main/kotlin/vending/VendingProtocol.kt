package vending

import mu.KotlinLogging
import java.net.ConnectException
import kotlin.experimental.and
import kotlin.experimental.inv
import kotlin.experimental.or

class VendingProtocol internal constructor() {
    private val logger = KotlinLogging.logger {}

    private lateinit var socket: VendingSocket

    private var timer: java.util.Timer? = null
    var status: Status = Status()

    fun connect(host: String, port: Int = 999) {
        try {
            socket = VendingSocket(host, port) //socket = new socket("192.168.232.2", 1024);
        } catch (e: ConnectException) {
            logger.error("Could not connect to coffee vending. $e")
            throw e
        }

        val socketTread = Thread(socket)
        socketTread.start()
        logic(initEvent, null)
    }

    fun disconnect() {
        socket.stop()
    }

    private fun parsePollData(pollData: ByteArray?): Int //TODO: Complete this.
    {
        return if (pollData != null) {
            if (pollData[2] == errorCode) errorCode.toInt() else 0
        } else {
            0
        }
    }

    fun syncIfNeed(data: ByteArray): Int {
        return if (data[2] == errorCode) {
            if (cmdSync() == 0) {
                YES
            } else noAnsw
        } else NO
    }

    private val YES = 16
    private val NO = -16
    private val noAnsw = -17
    val success = 32
    private val fail = -32
    private var lastMakeCmdStatus = fail
    private var pollStatus = fail
    private var syncRequestStatus = fail
    private var error = 0

    fun logic(source: Int, data: IntArray?): Int {
        if (source != isSellSucceedEvent && logicBusy) return -2
        if (error > 1) {
            error = 0
            logicBusy = false
            return -1
        }
        logicBusy = true
        when (source) {
            timerEvent -> {
                if (cmdPoll() == 0) {
                    val pollData: ByteArray? = socket.getAnswer()
                    if (parsePollData(pollData) == errorCode.toInt()) {
                        cmdSync() //todo check for results
                    }
                } else {
                    toLog("Reconnecting...")
                    socket.reconnect()
                    error++
                }
            }

            makeACoffeeEvent -> {
                toLog("makeACoffeeEvent")
                if (data != null) {
                    if (cmdSell(data[0], data[1]) == 0) {
                        toLog("cmdSellSended, answ = " + bytesToHex(socket.getAnswer()))
                        val answ: ByteArray? = socket.getAnswer()
                        if (answ != null) {
                            if (answ[2].toInt() == 0) //no errors
                            {
                                lastMakeCmdStatus = success //TODO УБРАТЬ!
                                val sellStatus = logic(isSellSucceedEvent, null) //TODO Check for process
                                return if (sellStatus == YES) {
                                    lastMakeCmdStatus = success
                                    logicBusy = false
                                    success //TODO change to 0
                                } else {
                                    logicBusy = false
                                    fail //TODO change to -7
                                }
                            } else {
                                toLog("cmdSellFail")
                                if (answ[2] == errorCode) {
                                    cmdSync()
                                    logicBusy = false
                                    logic(source, data)
                                }
                            }
                        } else {
                            logger.error("Could not get answer from socket")
                        }
                    } else {
                        toLog("No response to cmdSell. Retrying...")
                        socket.reconnect()
                        error++
                        logicBusy = false
                        logic(source, data)
                        error = 0
                        return if (lastMakeCmdStatus == success) 0 else -7
                    }
                } else {
                    logicBusy = false
                    logger.error("logic makeACoffee(): Be careful next time!")
                    return -3
                }
            }

            checkStatusEvent -> {}
            isSellSucceedEvent -> {
                toLog("isSellSucceedEvent")
                return if (cmdIsSellSucceed() == 0) {
                    val answ: ByteArray? = socket.getAnswer()
                    if (answ != null) {
                        if (answ[4].toInt() == 0x01) {
                            toLog("Sell succeed")
                            logicBusy = false
                            YES
                        } else {
                            toLog("Sell failed?")
                            logicBusy = false
                            NO
                        }
                    } else {

                        isSellSucceedEvent
                    }
                } else {
                    toLog("No response to cmdIsSellSucceed")
                    logicBusy = false
                    noAnsw
                }
            }

            pollRequest -> {
                run {
                    if (cmdPoll() == 0) {
                        pollData = socket.getAnswer()
                        if (pollData!![2] != errorCode) {
                            pollStatus = success
                            logicBusy = false
                            return success
                        } else {
                            logicBusy = false
                            if (logic(
                                    syncRequest, null
                                ) == success
                            ) {
                                logic(pollRequest, null)
                            }
                        }
                    } else {
                        toLog("No responce to cmdPoll. Retrying...")
                        socket.reconnect()
                        error++
                        logicBusy = false
                        logic(source, data)
                        error = 0
                        return if (pollStatus == success) success else fail /*
                    if(cmdSync() == 0)
                    {
                        if(syncIfNeed(socket.getAnswer())== noAnsw) {
                            logicBusy = false;
                            return fail;
                        }
                    }
                    else {
                        logicBusy = false;
                        return fail;
                    }
                    if(cmdPoll() == 0) {
                        logicBusy = false;
                        return success;
                    }
                    else {
                        logicBusy = false;
                        return fail;
                    }*/
                    }
                }
                run {
                    return if (cmdSync() == 0) {
                        logicBusy = false
                        syncRequestStatus = success
                        error = 0
                        success
                    } else {
                        socket.reconnect()
                        error++
                        logic(source, data)
                        error = 0
                        logicBusy = false
                        if (syncRequestStatus == fail) fail else success
                    }
                }
            }

            syncRequest -> {
                return if (cmdSync() == 0) {
                    logicBusy = false
                    syncRequestStatus = success
                    error = 0
                    success
                } else {
                    socket.reconnect()
                    error++
                    logic(source, data)
                    error = 0
                    logicBusy = false
                    if (syncRequestStatus == fail) fail else success
                }
            }

            initEvent -> {
                timer =
                    java.util.Timer() //                                                                    TODO Enable!
            }
        }
        logicBusy = false
        return 0
    }

    var logicBusy = false
    var pollData: ByteArray? = null

    inner class Status internal constructor() {
        var a = 0
        var b = 0
        val isCupIsStillHere: Boolean
            get() {
                if (pollData != null) {
                    if (pollData!![13].toInt() and 0x0C == 0x0C) {
                        return true
                    }
                }
                return false
            }
        val isVendingIdle: Boolean
            get() {
                if (pollData != null) {
                    if (pollData!![13].toInt() == 0x02) return true
                }
                return false
            }

        //сбой продажи
        val isSomethingHappend: Boolean
            get() {
                if (pollData != null) if (pollData!![13].toInt() == 0x0D) return true
                return false
            }
    }

    var poll: java.util.TimerTask = object : java.util.TimerTask() {
        override fun run() {
            logic(timerEvent, null)
        }
    }

    fun isCrcCorrect(data: ByteArray): Boolean {
        var crc = 0xAA
        var i = 1
        while (i < data[3] + 3 /*+ 1*/ && i < data.size) {
            if (data[i] < 0) crc += +0x100 //TODO check
            crc += data[i].toInt()
            i++
        }
        return crc == data[4].toInt()
    }

    fun recoverBadBytes(data: ByteArray): ByteArray {
        var reduce = 0
        for (datum in data) if (datum.toInt() == 0x28) reduce++
        val tmp: java.nio.ByteBuffer = java.nio.ByteBuffer.allocate(data.size - reduce)
        var i = 0
        while (i < data.size) {
            if (data[i].toInt() == 0x28) {
                i++
                tmp.put((data[i] and 0xF0.toByte() or (data[i].inv() and 0x0F)))
            } else tmp.put(data[i])
            i++
        }
        return tmp.array()
    }

    fun unformAnswer(data: ByteArray): ByteArray? {
        val badData = recoverBadBytes(data)
        return if (isCrcCorrect(badData)) badData else null
    }

    fun crc(data: ByteArray): Byte {
        var crc = 0xAA
        var i = 1
        while (i < data[3] + 3 + 1 && i < data.size) {
            if (data[i] < 0) crc += +0x100
            crc += data[i].toInt()
            i++
        }
        return (crc and 0xFF).toByte()
    }

    fun replaceBadBytes(data: ByteArray): ByteArray {
        var expand = 0
        for (i in 1 until data.size) if (data[i] == 0xD7.toByte() || data[i] == 0x28.toByte()) expand++
        if (expand == 0) return data
        val tmp: java.nio.ByteBuffer = java.nio.ByteBuffer.allocate(data.size + expand)
        tmp.put(0xD7.toByte())
        for (i in 1 until data.size) {
            if (data[i] == 0xD7.toByte() || data[i] == 0x28.toByte()) {
                tmp.put(0x28.toByte())
                tmp.put(data[i] and 0xF0.toByte() or data[i].inv() and 0x0F.toByte())
            } else tmp.put(data[i])
        }
        return tmp.array()
    }

    var currentID = 0
    fun formRequest(data: ByteArray): Int {
        var data = data
        data[0] = 0xD7.toByte()
        data[1] = (currentID and 0xFF).toByte()
        data[4 + data[3]] = crc(data)
        data = replaceBadBytes(data)
        toLog("SENDING: " + bytesToHex(data))
        val status: Int = socket.send(data)
        if (status == 0) if (currentID == 255) currentID = 0 else currentID++
        toLog("Received: " + bytesToHex(socket.getAnswer()))
        return status
    }

    fun cmdSync(): Int {
        val data = ByteArray(5)
        data[2] = 0x00
        data[3] = 0
        return formRequest(data)
    }

    fun cmdPoll(): Int {
        val data = ByteArray(5)
        data[2] = 0x01
        data[3] = 0
        return formRequest(data)
    }

    fun cmdTest(): Int {
        val data = ByteArray(15)
        data[2] = 100
        data[3] = 10
        data[4] = 0xD7.toByte()
        data[5] = 0xD7.toByte()
        data[6] = 0
        data[7] = 0x28.toByte()
        return formRequest(data)
    }

    fun cmdSell(coffeeType: Int, sugar: Int): Int {
        val data = ByteArray(3 + 5)
        data[2] = 0x04
        data[3] = 3
        data[4] = 0x20 //код аппарата ()
        data[5] = (coffeeType and 0xFF).toByte()
        data[6] = (0x80 + sugar).toByte()
        return formRequest(data)
    }

    fun cmdIsSellSucceed(): Int {
        val data = ByteArray(5)
        data[2] = 0x05
        data[3] = 0
        return formRequest(data)
    }

    val socketLog: String
        get() = bytesToHex(socket.getAnswer())

    /*public void cmdSell(int coffeeType, int sugar) {
          if(cmdSell(13, 0) == 0)
          {
              toLog(bytesToHex(socket.getAnswer()));
          }
      }

      public void onCmdSync(View view) {
          if(cmdSync() == 0)
          {
              toLog(bytesToHex(socket.getAnswer()));
          }
      }

      public void onCmdPoll(View view) {
          if(cmdPoll() == 0)
          {
              toLog(bytesToHex(socket.getAnswer()));
          }
      }

      public void onMakeACoffee(View view) {
          if(makeACoffee(Integer.valueOf(coffeeType.getText().toString()), 0) == success)
              toLog("Make a coffee success!");
          else
              toLog("Sell failed!" + bytesToHex(socket.getAnswer()));
      }*/
    fun makeACoffee(type: Int, sugar: Int): Int {
        lastMakeCmdStatus = fail
        return logic(makeACoffeeEvent, intArrayOf(type, sugar))
    }

    var logData: String? = null

    fun toLog(text: String) {
        logData += java.util.Calendar.getInstance().getTime().toString() + " " + text + "\n"
    }

    @JvmName("getLogData1")
    fun getLogData(): String? {
        val toSend = logData
        logData = null
        return toSend
    }

    fun readPoll(): Int {
        return logic(pollRequest, null)
    }

    companion object {
        const val errorCode = 0xFE.toByte()

        /*void retry(Function <Void, Void> t)
    {
        t.apply();

    }*/
        const val initEvent = 0
        const val timerEvent = 1
        const val makeACoffeeEvent = 2
        const val checkStatusEvent = 3 //needs for get function
        const val isSellSucceedEvent = 4
        const val pollRequest = 5
        const val syncRequest = 6
        private val hexArray: CharArray = "0123456789ABCDEF".toCharArray()
        fun bytesToHex(bytes: ByteArray?): String {
            if (bytes == null) return "nothing"
            val hexChars = CharArray(bytes.size * 3)
            var j = 0
            while (j < bytes.size && j < bytes[3] + 5) {
                val v = bytes[j].toInt() and 0xFF
                hexChars[j * 3] = hexArray[v ushr 4]
                hexChars[j * 3 + 1] = hexArray[v and 0x0F]
                hexChars[j * 3 + 2] = ' '
                j++
            }
            return String(hexChars)
        }
    }
}
