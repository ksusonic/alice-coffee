package vending

import kotlinx.coroutines.runBlocking
import mu.KotlinLogging
import java.io.IOException
import java.io.InputStream
import java.io.OutputStream
import java.net.Socket

class VendingSocket internal constructor(private var host: String, private var port: Int) : Runnable, AbstractSocket {
    private val logger = KotlinLogging.logger {}

    private var socket: Socket = Socket(host, port)
    private var socketError = "SocketError"
    private var testOut: OutputStream? = null
    private var testIn: InputStream? = null
    private var dataToSend: ByteArray? = null
    private var answer: ByteArray? = null
    private val timeoutDefault = 700 //change to 200 for real machine
    private var timeout = timeoutDefault

    private var running = true
    private var read = false
    private var reconnect = false
    private var waitForRead = 0

    override suspend fun connect() {
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

    override suspend fun reconnect() {
        disconnect()
        connect()
    }

    override fun run() {
        runBlocking {
            connect()
        }
        var startTime: Long
        while (running) {
            startTime = System.currentTimeMillis()
            if (reconnect) {
                runBlocking {
                    reconnect()
                }
            }
            if (dataToSend != null) {
                runBlocking {
                    send(dataToSend!!)
                }
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

    override suspend fun send(data: ByteArray): Int {
        write(data)
        dataToSend = null
        return 0
    }

    override fun getAnswer(): ByteArray? {
        return answer
    }
}