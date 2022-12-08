import org.apache.logging.log4j.LogManager
import java.io.IOException
import java.io.InputStream
import java.io.OutputStream
import java.net.Socket

class VendingSocket internal constructor(private var host: String, private var port: Int) : Runnable {
    private val logger = LogManager.getLogger()

    private var socket: Socket = Socket(host, port)
    private var socketError = "SocketError"
    private var testOut: OutputStream? = null
    private var testIn: InputStream? = null
    private var dataToSend: ByteArray? = null
    private var answer: ByteArray? = null
    private val timeoutDefault = 700 //change to 200 for real machine
    private val timeoutLong = 10000 //change to 1000 for real machine
    private var timeout = timeoutDefault

    private var running = true
    private var read = false
    private var reconnect = false
    private var waitForRead = 0

    fun connect() {
        try {
            socket = Socket(host, port)
            testIn = socket.getInputStream()
            testOut = socket.getOutputStream()
        } catch (e: IOException) {
            logger.error(socketError, e)
        } catch (e: NullPointerException) {
            logger.error(socketError, e)
        }
    }

    private fun disconnect() {
        try {
            socket.close()
        } catch (e: IOException) {
            e.printStackTrace()
        } catch (e: NullPointerException) {
            e.printStackTrace()
        }
    }

    private fun Reconnect() {
        disconnect()
        connect()
    }

    fun reconnect() {
        reconnect = true
        try {
            Thread.sleep(500)
        } catch (e: InterruptedException) {
            e.printStackTrace()
        }
    }

    fun stop() {
        running = false
    }

    override fun run() {
        connect()
        var startTime: Long
        while (running) {
            startTime = System.currentTimeMillis()
            if (reconnect) {
                reconnect = false
                Reconnect()
            }
            if (dataToSend != null) {
                send(dataToSend!!)
            }
            if (read) {
                read = false
                answer = read()
                waitForRead = if (answer == null) -1 else 0
            }
            if (500 - (System.currentTimeMillis() - startTime) > 0) try {
                Thread.sleep(500 - (System.currentTimeMillis() - startTime)) //отдыхай
            } catch (_: InterruptedException) {
            }
        }
        try {
            socket.close()
        } catch (_: IOException) {
        } catch (_: NullPointerException) {
        }
    }

    private fun write(Data: ByteArray): Int {
        try {
            testOut!!.write(Data)
            testOut!!.flush()
        } catch (e: IOException) {
            logger.error(socketError, "Socket write error.")
            return -1
        } catch (e: NullPointerException) {
            logger.error(socketError, "Socket write error.")
            return -1
        }
        return 0
    }

    private fun read(): ByteArray? {
        val data = ByteArray(1023)
        return try {
            socket.soTimeout = timeout
            val result = testIn!!.read(data)
            if (result > 0) data else null
        } catch (e: IOException) {
            logger.error(socketError, "Socket read error.")
            null
        } catch (e: NullPointerException) {
            logger.error(socketError, "Socket read error.")
            null
        }
    }

    private fun send(data: ByteArray) {
        write(data)
        dataToSend = null
    }

    fun send(data: ByteArray?): Int {
        dataToSend = data
        read = true
        waitForRead = -2
        while (waitForRead == -2) {
            try {
                Thread.sleep(100)
            } catch (e: InterruptedException) {
                e.printStackTrace()
            }
        }
        return waitForRead
    }

    fun getAnswer(): ByteArray? {
        return answer
    }
}