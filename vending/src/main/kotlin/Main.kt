import kotlinx.coroutines.runBlocking
import mu.KotlinLogging
import server.ServerController

private val logger = KotlinLogging.logger {}

fun main() {
    logger.info("Starting Alice coffee app!")

    val vending = vending.VendingProtocol("127.0.0.1", 8081)
    val server = ServerController("ws://127.0.0.1:8080/ws", vending)

    vending.connect()
    runBlocking {
        server.connect()
    }
}
