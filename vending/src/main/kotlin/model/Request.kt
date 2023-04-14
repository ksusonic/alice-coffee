package model

import kotlinx.serialization.Serializable

@Serializable
data class Request(
    // reqid
    val id: String,
    // type of coffee
    val type: String,
    val sugar: Int
)
