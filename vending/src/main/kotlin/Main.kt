import kotlinx.coroutines.runBlocking
import mu.KotlinLogging
import server.ServerController
import kotlin.system.exitProcess

private val logger = KotlinLogging.logger {}

fun main() {
    logger.info("Starting Alice coffee app!")

    val vending = vending.VendingProtocol("127.0.0.1", 8081)
    val server = ServerController("ws://127.0.0.1:8080/ws", vending)

    logger.info("starting cloud-client")
    runBlocking {
        while (true) {
            try {
                server.serve()
            } catch (e: kotlinx.coroutines.channels.ClosedReceiveChannelException) {
                logger.error("caught ClosedReceiveChannelException: $e")
            } catch (e: java.net.ConnectException) {
                logger.error("cannot reconnect to vending. Shutting down.")
                exitProcess(1)
            }
        }
    }
}
