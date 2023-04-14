package server

import io.ktor.client.*
import io.ktor.client.plugins.websocket.*
import io.ktor.websocket.*
import vending.VendingProtocol
import java.util.*


class ServerController(
    private val urlString: String,
    private val vending: VendingProtocol
) {
    private val client = HttpClient {
        install(WebSockets)
    }

    suspend fun connect() = client.webSocket(urlString = urlString) {
        while (true) {
            val othersMessage = incoming.receive() as? Frame.Text
            println(othersMessage?.readText())
            val myMessage = Scanner(System.`in`).next()
            if (myMessage != null) {
                send(myMessage)
            }
        }
    }

}