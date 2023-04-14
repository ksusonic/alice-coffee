package com.example

import io.ktor.network.selector.*
import io.ktor.network.sockets.*
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.runBlocking

fun main() {
    runBlocking {
        val selectorManager = SelectorManager(Dispatchers.IO)
        val serverSocket = aSocket(selectorManager).tcp().bind("127.0.0.1", 8081)
        println("Server is listening at ${serverSocket.localAddress}")
        while (true) {
            val socket = serverSocket.accept()
            println("Accepted from ${socket.remoteAddress}")

            val receiveChannel = socket.openReadChannel()
            try {
                while (true) {
                    val byte = receiveChannel.readByte()
                    println("got: $byte")
                }
            } catch (e: Throwable) {
                println("got error: $e")
                socket.close()
            }
        }
    }
}