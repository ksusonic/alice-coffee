import mu.KotlinLogging
import queue.MessageHandler

private val logger = KotlinLogging.logger {}

suspend fun main() {
    logger.info("Starting Alice coffee app!")

    val mq = MessageHandler()

//    val vending = vending.VendingProtocol()
    // vending.connect()

    while (true) {
        val rawCommand = mq.receiveAndDeleteMessage()
        logger.debug("got command: $rawCommand")
    }
}
