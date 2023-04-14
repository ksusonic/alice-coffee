package server

import io.ktor.client.*
import io.ktor.client.plugins.*
import io.ktor.client.plugins.websocket.*
import io.ktor.serialization.kotlinx.*
import kotlinx.serialization.json.Json
import model.Request
import mu.KotlinLogging
import vending.VendingProtocol


class ServerController(
    private val urlString: String,
    private val vending: VendingProtocol
) {
    private val logger = KotlinLogging.logger {}

    private val client = HttpClient {
        install(WebSockets) {
            contentConverter = KotlinxWebsocketSerializationConverter(Json)
        }
        install(HttpRequestRetry) {
            retryOnServerErrors(maxRetries = 5)
            delayMillis { retry ->
                retry * 3000L
            }
            exponentialDelay() // 3 6 18
        }
    }

    suspend fun serve() = client.webSocket(urlString) {
        logger.info("connected to ws")
        val request = receiveDeserialized<Request>()

        logger.info("got request ${request.id}")
        if (request.type.isNotEmpty()) {
            vending.makeACoffee(2, request.sugar)
        }
    }

}