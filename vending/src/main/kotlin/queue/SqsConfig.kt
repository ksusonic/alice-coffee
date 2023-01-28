package queue

import aws.sdk.kotlin.runtime.endpoint.AwsEndpoint
import aws.sdk.kotlin.runtime.endpoint.AwsEndpointResolver
import aws.sdk.kotlin.runtime.endpoint.CredentialScope
import aws.smithy.kotlin.runtime.auth.awscredentials.Credentials
import aws.smithy.kotlin.runtime.auth.awscredentials.CredentialsProvider
import java.lang.System.getenv

const val SQS_REGION = "ru-central1"
const val YA_CLOUD_API_URL = "https://message-queue.api.cloud.yandex.net"

class SqsConfigLoader : CredentialsProvider {
    override suspend fun getCredentials(): Credentials = Credentials(
        accessKeyId = getenv("AWS_ACCESS_KEY_ID"),
        secretAccessKey = getenv("AWS_SECRET_ACCESS_KEY")
    )
}

class SqsEndpointResolver(queueApiUrl: String = YA_CLOUD_API_URL) : AwsEndpointResolver {
    private val apiUrl = queueApiUrl
    override suspend fun resolve(service: String, region: String) = AwsEndpoint(
        apiUrl, CredentialScope(region = SQS_REGION)
    )
}
