package main.emulation

import vending.VendingProtocol
import kotlin.test.assertEquals
import kotlin.test.assertTrue


class Emulator {
    val syncCode: Byte = 0x00
    val SellCommand: Byte = 0x04
    val IsSellSucceedCommand: Byte = 0x05

    val VendingCode: Byte = 0x20

    var isSynced: Boolean = false
    var lastSellSuccessed: Boolean = false


    data class Response(
        val message: String,
        val response: ByteArray = byteArrayOf(0, 0, 0, 0, 0, 0, 0, 0)
    )

    fun handle(data: ByteArray): Response {
        assertEquals(8, data.size, "expected package of 8 bytes")
        assertEquals(0xD7.toByte(), data[0], "every request must start with 0xD7")
        assertTrue(data[1] >= 0, "request number must me >= 0.")

        if (data[2] == syncCode) {
            isSynced = true
            println("synced. Ok")
            return Response("synced. Ok")
        }
        if (!isSynced) {
            throw RuntimeException("ERROR: NOT SYNCED")
        }
        when (data[2]) {
            SellCommand -> {
                println("handling make coffee command")
                if (data[4] != VendingCode) {
                    println("vending code is not $VendingCode. Got: ${data[4]}")
                }
                val type = VendingProtocol.CoffeeType.values().find {
                    it.value == data[5].toInt()
                } ?: VendingProtocol.CoffeeType.UNKNOWN
                lastSellSuccessed = true
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

            IsSellSucceedCommand -> {
                println("handling is last sell successed command")
                return Response(
                    lastSellSuccessed.toString(), byteArrayOf(
                        data[0],
                        data[1],
                        0, // no errors
                        0,
                        if (lastSellSuccessed) {
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
//        assertEquals(VendingCode, data[2], "code of vending (data[2]) is not correct")

        return Response("ok")
    }
}