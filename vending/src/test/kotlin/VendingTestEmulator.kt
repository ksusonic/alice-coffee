package vending.emulator

import io.ktor.network.selector.*
import io.ktor.network.sockets.*
import io.ktor.utils.io.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.runBlocking


suspend fun read8Bytes(input: ByteReadChannel): ByteArray? {
    val buffer = ByteArray(8) // create a buffer to store 8 received bytes
    val totalBytesRead = 0 // initialize total bytes read to 0

    // read exactly 8 bytes
    val bytesRead = input.readAvailable(
        buffer,
        totalBytesRead,
        8
    ) // read available bytes from the input channel into the buffer
    if (bytesRead == -1) return null // if the channel is closed, return null
    return buffer
}


fun main() {
    runBlocking {
        val selectorManager = SelectorManager(Dispatchers.IO)
        val serverSocket = aSocket(selectorManager).tcp().bind("127.0.0.1", 8081)
        val emulator = Emulator()

        println("Server is listening at ${serverSocket.localAddress}")
        while (true) {
            var socket = serverSocket.accept()
            println("Accepted from ${socket.remoteAddress}")

            val receiveChannel = socket.openReadChannel()
            val writeChannel = socket.openWriteChannel(autoFlush = false)

            while (true) {
                if (socket.isClosed) {
                    socket = serverSocket.accept()
                }
                var message: ByteArray
                try {
                    message = read8Bytes(receiveChannel) ?: throw RuntimeException("message is null")
                } catch (e: Exception) {
                    println("lost connection")
                    break
                }

                println("got message: $message")
                try {
                    val result = emulator.handle(message)
                    println("Successful result: $result")
                    for (b in result.response) {
                        writeChannel.writeByte(b)
                    }
                    writeChannel.flush()
                    println("response with: ${result.response}")
                } catch (e: Exception) {
                    println("Could not emulate: ${e.message}")
                }
            }
        }
    }
}
