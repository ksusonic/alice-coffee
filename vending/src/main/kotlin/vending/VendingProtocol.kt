package vending

import kotlinx.coroutines.runBlocking
import mu.KotlinLogging
import kotlin.experimental.and
import kotlin.experimental.inv
import kotlin.experimental.or

interface AbstractSocket {
    suspend fun reconnect()
    fun getAnswer(): ByteArray?
    suspend fun send(data: ByteArray): Int
    suspend fun connect()
}

class VendingProtocol internal constructor(private val socket: AbstractSocket) {

    enum class CoffeeType(val value: Int, val readable: String, val code: String) {
        BLACK(0, "черный кофе", "black"),
        ESPRESSO(1, "эспрессо", "espresso"),
        AMERICANO(2, "американо", "americano"),
        LATTE(3, "латте", "latte"),
        CAPPUCCINO(4, "капучино", "cappuccino"),
        DOUBLE(8, "двойной", "double"),
        TEA(126, "чай", "tea"),
        UNKNOWN(128, "UNKNOWN", "unknown")
    }

    private val logger = KotlinLogging.logger {}

    private var timer: java.util.Timer? = null
//    var status: Status = Status() TODO: status

    private fun parsePollData(pollData: ByteArray?): Int //TODO: Complete this.
    {
        return if (pollData != null) {
            if (pollData[2] == errorCode) errorCode.toInt() else 0
        } else {
            0
        }
    }

    private val yes = 16
    private val no = -16
    private val noAnswer = -17
    private val success = 32
    private val fail = -32
    private var lastMakeCmdStatus = fail
    private var pollStatus = fail
    private var syncRequestStatus = fail
    private var error = 0

    suspend fun logic(source: Int, data: IntArray?): Int {
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
                    logger.info("Reconnecting...")
                    socket.reconnect()
                    error++
                }
            }

            makeACoffeeEvent -> {
                logger.info("makeACoffeeEvent")
                if (data != null) {
                    val sellResult = cmdSell(data[0], data[1])
                    if (sellResult == 0) {
                        logger.debug("cmdSellSended, answ = " + bytesToHex(socket.getAnswer()))
                        val answ: ByteArray? = socket.getAnswer()
                        if (answ != null) {
                            if (answ[2].toInt() == 0) //no errors
                            {
                                lastMakeCmdStatus = success //TODO УБРАТЬ!
                                val sellStatus = logic(isSellSucceedEvent, null) //TODO Check for process
                                return if (sellStatus == yes) {
                                    lastMakeCmdStatus = success
                                    logicBusy = false
                                    success //TODO change to 0
                                } else {
                                    logicBusy = false
                                    fail //TODO change to -7
                                }
                            } else {
                                logger.error("cmdSellFail")
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
                        logger.error("No response to cmdSell. Retrying...")
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
                logger.debug("isSellSucceedEvent")
                return if (cmdIsSellSucceed() == 0) {
                    val answ: ByteArray? = socket.getAnswer()
                    if (answ != null) {
                        if (answ[4].toInt() == 0x01) {
                            logger.info("Sell succeed")
                            logicBusy = false
                            yes
                        } else {
                            logger.info("Sell failed?")
                            logicBusy = false
                            no
                        }
                    } else {

                        isSellSucceedEvent
                    }
                } else {
                    logger.error("No response to cmdIsSellSucceed")
                    logicBusy = false
                    noAnswer
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
                        logger.error("No responce to cmdPoll. Retrying...")
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

    private var logicBusy = false
    var pollData: ByteArray? = null

    inner class Status internal constructor() {
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
            runBlocking {
                logic(timerEvent, null)
            }
        }
    }

    private fun isCrcCorrect(data: ByteArray): Boolean {
        var crc = 0xAA
        var i = 1
        while (i < data[3] + 3 /*+ 1*/ && i < data.size) {
            if (data[i] < 0) crc += +0x100 //TODO check
            crc += data[i].toInt()
            i++
        }
        return crc == data[4].toInt()
    }

    private fun recoverBadBytes(data: ByteArray): ByteArray {
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

    private fun crc(data: ByteArray): Byte {
        var crc = 0xAA
        var i = 1
        while (i < data[3] + 3 + 1 && i < data.size) {
            if (data[i] < 0) crc += +0x100
            crc += data[i].toInt()
            i++
        }
        return (crc and 0xFF).toByte()
    }

    private fun replaceBadBytes(data: ByteArray): ByteArray {
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

    private var currentID = 0
    private suspend fun formRequest(data: ByteArray): Int {
        var data = data
        data[0] = 0xD7.toByte()
        data[1] = (currentID and 0xFF).toByte()
        data[4 + data[3]] = crc(data)
        data = replaceBadBytes(data)
        logger.debug("SENDING: " + bytesToHex(data))
        val status: Int = socket.send(data)
        if (status == 0) if (currentID == 255) currentID = 0 else currentID++
        logger.debug("Received: " + bytesToHex(socket.getAnswer()))
        return status
    }

    suspend fun cmdSync(): Int {
        val data = ByteArray(5)
        data[2] = 0x00
        data[3] = 0
        return formRequest(data)
    }

    private suspend fun cmdPoll(): Int {
        val data = ByteArray(5)
        data[2] = 0x01
        data[3] = 0
        return formRequest(data)
    }

    suspend fun cmdTest(): Int {
        val data = ByteArray(15)
        data[2] = 100
        data[3] = 10
        data[4] = 0xD7.toByte()
        data[5] = 0xD7.toByte()
        data[6] = 0
        data[7] = 0x28.toByte()
        return formRequest(data)
    }

    private suspend fun cmdSell(coffeeType: Int, sugar: Int): Int {
        val data = ByteArray(3 + 5)
        data[2] = 0x04
        data[3] = 3
        data[4] = 0x20 //код аппарата ()
        data[5] = (coffeeType and 0xFF).toByte()
        data[6] = (0x80 + sugar).toByte()
        return formRequest(data)
    }

    private suspend fun cmdIsSellSucceed(): Int {
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
    suspend fun makeACoffee(type: Int, sugar: Int): Int {
        lastMakeCmdStatus = fail
        return logic(makeACoffeeEvent, intArrayOf(type, sugar))
    }

    suspend fun readPoll(): Int {
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
