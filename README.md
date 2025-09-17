# UUID

The RFC4122 specifies seven variants of UUID ([New UUID Formats](https://uuid6.github.io/uuid6-ietf-draft/)).:
* **Version 1:** Time-based UUID.
* **Version 3:** MD5 hash of some data.
* **Version 4:** Random data.
* **Version 5:** SHA1 hash of some data.
* **Version 6:** Timestamp and monotonic counter.
* **Version 7:** Unix timestamp.
* **Version 8:** user-defined data.

## Version 1: Time-based UUID
UUID (v1) creates a unique identifier by merging the MAC address with the current timestamp, resulting in a 128-bit value displayed as a 16-byte sequence. It’s often formatted as "8-4-4-4-12" for readability in various displays and records.

Here's a sample representation in hexadecimal form:

```
e29b44a0-4d2c-11eb-b378-0242ac130004
```

Storage format:
* **e29b44a0** (time_low, **32 bits**) => `11100010100110110100010010100000`
* **4d2c** (time_mid, **16 bits**) => `0100110100101100`
* **11eb** (time_hi_and_version, **16 bits**) => `0001 000111101011`
  * **4 bits** for version => `0001`
  * **12 bits** for time_hi => `000111101011`
* **b378** (clock_seq, **16 bits**) => `10 11001101111000`
  * **2-3 bits** for the variant (in this example, it starts with **10**, which aligns with RFC 4122) => `10`
  * **14 bits** for clock sequence => `11001101111000`
* **0242ac130004** (node, **48 bits**) => `000000100100001010101100000100110000000000100`

### Timestamp
[RFC 4122: Timestamp](https://www.rfc-editor.org/rfc/rfc4122#section-4.1.4) specifies timestamp as *Count of 100-nanosecond intervals since 00:00:00.00, 15 October 1582*. 
* **Current timestamp**: Current UTC time in 100-nanoseconds interval 
* **Offset**: 100-nanosecond intervals between the UUID epoch (15 October 1582) and the Unix epoch (1 January 1970).

### Clock Sequence
The purpose of the [RFC 4122: Clock Sequence](https://www.rfc-editor.org/rfc/rfc4122#section-4.1.5) is to handle cases where the system clock goes backward or to ensure uniqueness if the same UUID is generated more than once in the same *100ns* interval (timestamp precision).

Clock sequence is increment in the following scenarios:
- Upon initialization
- On update

### MAC Address
- *r.Intn(1 << 14)* generates 14bit number
- *Increment* ensures two things:
  * The value of **ClockSequence** is incremented.
  * The resulting value doesn't exceed the 14-bit limit.
* **0x3FFF** is a hexadecimal representation of a number with the last 14 bits set to 1 and the rest set to 0. 

## Version 3: MD5 hash of some data

In the representation below, each 'x' is a placeholder for a hexadecimal character, and each hexadecimal character corresponds to 4 bits:
| time_low | time_mid | time_hi_ver | clock_seq_hi | clock_seq_low | node         |
|----------|----------|-------------|--------------|---------------|--------------|
| xxxxxxxx | xxxx     | 1xxx        | 1xxx         | 8xxx          | xxxxxxxxxxxx |

### Clock Sequence

The purpose of the [RFC 4122: Clock Sequence](https://www.rfc-editor.org/rfc/rfc4122#section-4.1.5) is to handle cases where the system clock goes backward or to ensure uniqueness if the same UUID is generated more than once in the same *100ns* interval (timestamp precision).

The clock sequence is a 14-bit value. Value is split in two parts:
* **clock_seq_hi_and_reserved** (8 bits): Upper 6 bits of the clock sequence combined with the 2-bit variant.
* **clock_seq_low** (8 bits): Lower 8 bits of the clock sequence.

Resources:
- [Clock Sequence Generation? · Issue \#41 · uuid6/uuid6-ietf-draft](https://github.com/uuid6/uuid6-ietf-draft/issues/41)

Here's a breakdown of the data length in bits for each component of a UUIDv1:
| Component      | Length (bits) |
|----------------|---------------|
| Timestamp      | 60            |
| Version        | 4             |
| Clock Sequence | 14            |
| Node           | 48            |
* **time_low**: The first 32 bits of the timestamp.
* **time_mid**: The next 16 bits of the timestamp.
* **version**: 4 bits, where the version number is stored. For version 1, this will be **1**.
* **time_high**: The next 12 bits of the timestamp.
* **clock_seq_hi_res** and **clock_seq_low**: Together, these form the 14-bit clock sequence.
* **node**: 48 bits representing the node, usually the MAC address of the host.
## Version 1

* **Version 1:** UUIDs using a timestamp and monotonic counter.
* **Version 3:** UUIDs based on the MD5 hash of some data.
* **Version 4:** UUIDs with random data.
* **Version 5:** UUIDs based on the SHA1 hash of some data.
* **Version 6:** UUIDs using a timestamp and monotonic counter.
* **Version 7:** UUIDs using a Unix timestamp.
* **Version 8:** UUIDs using user-defined data.
