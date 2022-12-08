import org.apache.logging.log4j.LogManager
import org.apache.logging.log4j.Logger

val logger: Logger = LogManager.getLogger()

fun main(args: Array<String>) {
    logger.debug("Starting app: $args")
    println("Hello World!")
    VendingProtocol()
}
