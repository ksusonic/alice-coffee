package queue

suspend fun main() {
    println("Starting debug app for message queue")
    while (true) {
        MessageHandler().receiveAndDeleteMessage()
    }
}