package vending.emulator

import vending.VendingProtocol
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class Emulator {
    private val syncCode: Byte = 0x00
    private val sellCommand: Byte = 0x04
    private val isSellSucceedCommand: Byte = 0x05

    private val vendingCode: Byte = 0x20

    private var isSynced: Boolean = false
    private var lastSellSuccess: Boolean = false


    data class Response(
        val message: String,
        val response: ByteArray = byteArrayOf(0, 0, 0, 0, 0, 0, 0, 0)
    ) {
        override fun equals(other: Any?): Boolean {
            if (this === other) return true
            if (javaClass != other?.javaClass) return false

            other as Response

            return response.contentEquals(other.response)
        }

        override fun hashCode(): Int = response.contentHashCode()
    }

    fun handle(data: ByteArray): Response {
        assertEquals(8, data.size, "expected package of 8 bytes")
        assertEquals(0xD7.toByte(), data[0], "every request must start with 0xD7")
        assertTrue(data[1] >= 0, "request number must me >= 0.")

        if (data[2] == syncCode) {
            isSynced = true
            println("synced. Ok")
            return Response("synced. Ok")
        }
        when (data[2]) {
            sellCommand -> {
                println("handling make coffee command")
                if (!isSynced) {
                    throw RuntimeException("ERROR: NOT SYNCED")
                }
                isSynced = false

                if (data[4] != vendingCode) {
                    println("vending code is not $vendingCode. Got: ${data[4]}")
                }
                val type = VendingProtocol.CoffeeType.values().find {
                    it.value == data[5].toInt()
                } ?: VendingProtocol.CoffeeType.UNKNOWN
                lastSellSuccess = true
                val readable = "making a drink with type: ${type.readable} with ${data[6].toInt()} sugar"
                println(readable)
                return Response(
                    readable, byteArrayOf(
                        data[0],
                        data[1],
                        0, // no errors
                        0,
                        0x01,
                        0,
                        0,
                        0
                    )
                )
            }

            isSellSucceedCommand -> {
                println("got check for is last sell success command == $lastSellSuccess")
                return Response(
                    lastSellSuccess.toString(), byteArrayOf(
                        data[0],
                        data[1],
                        0, // no errors
                        0,
                        if (lastSellSuccess) {
                            0x01
                        } else {
                            0x00
                        },
                        0,
                        0,
                        0
                    )
                )
            }
        }

        return Response("ok")
    }
}