package main

import kotlinx.coroutines.runBlocking
import mu.KotlinLogging
import org.slf4j.event.Level
import server.ServerController
import vending.AbstractSocket
import vending.EmulationSocket
import vending.VendingSocket
import java.net.ConnectException
import kotlin.system.exitProcess

private val logger = KotlinLogging.logger {}

const val DefaultAddress = "127.0.0.1"
const val DefaultPort = 8081

suspend fun main(args: Array<String>) {
    logger.atLevel(Level.DEBUG)
    logger.info("welcome to alice coffee app!")

    val config = parseConfig(args)
    logger.debug(config.toString())

    val socket: AbstractSocket = if (!config.emulationMode) {
        VendingSocket(getAddress(), getPort())
    } else {
        EmulationSocket(getAddress(), getPort())
    }

    runBlocking {
        try {
            socket.connect()
        } catch (e: ConnectException) {
            logger.error("failed to connect to vending socket: ${e.message}")
        }
    }

    val vending = vending.VendingProtocol(socket)
    val server = ServerController("ws://127.0.0.1:8080/ws", vending)

    logger.info("starting cloud-client")
    runBlocking {
        while (true) {
            try {
                server.serve()
            } catch (e: kotlinx.coroutines.channels.ClosedReceiveChannelException) {
                logger.error("caught ClosedReceiveChannelException: $e")
            } catch (e: ConnectException) {
                logger.error("cannot reconnect to vending. Shutting down.")
                exitProcess(1)
            }
        }
    }
}

fun getAddress(): String {
    return DefaultAddress
}

fun getPort(): Int {
    return DefaultPort
}