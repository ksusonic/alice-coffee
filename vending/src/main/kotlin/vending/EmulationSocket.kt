package vending

import io.ktor.network.selector.*
import io.ktor.network.sockets.*
import io.ktor.utils.io.*
import kotlinx.coroutines.Dispatchers
import mu.KotlinLogging

class EmulationSocket(private val address: String, private val port: Int) : AbstractSocket {
    private val logger = KotlinLogging.logger {}

    private val selectorManager = SelectorManager(Dispatchers.IO)
    private val socketBuilder = aSocket(selectorManager).tcp()

    private lateinit var socket: Socket
    private lateinit var recieveChannel: ByteReadChannel
    private lateinit var sendChannel: ByteWriteChannel

    private var answer: ByteArray? = null

    override suspend fun connect() {
        socket = socketBuilder.connect(address, port)
        recieveChannel = socket.openReadChannel()
        sendChannel = socket.openWriteChannel(autoFlush = false)
    }

    override suspend fun reconnect() {
        socket.close()
        connect()
    }

    override fun getAnswer(): ByteArray? {
        return answer
    }

    override suspend fun send(data: ByteArray): Int {
        sendChannel.writeAvailable(data)
        sendChannel.flush()


        val recieveBuffer = ByteArray(8)
        val recieve = recieveChannel.readAvailable(
            recieveBuffer,
            0,
            8
        )
        if (recieve == -1) {
            logger.error("no answer from vending")
            return 1
        }

        logger.debug("response: {}", recieveBuffer)
        answer = recieveBuffer
        return recieveBuffer[2].toInt()
    }
}