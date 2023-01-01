package queue

import aws.sdk.kotlin.services.sqs.SqsClient
import aws.sdk.kotlin.services.sqs.model.DeleteMessageRequest
import aws.sdk.kotlin.services.sqs.model.Message
import aws.sdk.kotlin.services.sqs.model.ReceiveMessageRequest
import aws.sdk.kotlin.services.sqs.model.ReceiveMessageResponse
import mu.KotlinLogging

const val DEFAULT_SQS_URL =
    "https://message-queue.api.cloud.yandex.net/b1g9hok9aibdsjgn3lge/dj600000000al7ui050h/coffee-to-make"

class MessageHandler(queueUrl: String = DEFAULT_SQS_URL, regionVal: String = SQS_REGION) {
    private val logger = KotlinLogging.logger {}
    private val queueUrlVar = queueUrl

    private val client = SqsClient {
        credentialsProvider = SqsConfigLoader(); region = regionVal; endpointResolver = SqsEndpointResolver()
    }

    suspend fun receiveAndDeleteMessage(): String {
        val receiveMessageRequest = ReceiveMessageRequest {
            queueUrl = queueUrlVar
            maxNumberOfMessages = 1
        }

        var message: Message? = null
        while (message == null) {
            val response: ReceiveMessageResponse = client.receiveMessage(receiveMessageRequest)
            message = response.messages?.first()
            if (message != null) {
                logger.info("Got message id=${message.messageId} ${message.receiptHandle} - ${message.body}}")
            } else {
                logger.debug("Got no-messages response $response")
            }
        }
        client.deleteMessage(DeleteMessageRequest {
            queueUrl = queueUrlVar
            receiptHandle = message.receiptHandle
        })
        // fact of delete is not checked
        logger.info("Deleted message ${message.receiptHandle}")
        return message.body.orEmpty()
    }
}
