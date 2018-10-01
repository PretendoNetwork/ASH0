/*
 * Based off https://github.com/trapexit/wiiqt/blob/master/WiiQt/ash.cpp
 * This version doesn't use external libs and is 64 bit
 * Conversion by Nybbit
 */

#if defined(_MSC_VER)
#define BYTESWAP16 _byteswap_ushort
#define BYTESWAP32 _byteswap_ulong
#elif defined (__GNUC__)
#define BYTESWAP16 __builtin_bswap16
#define BYTESWAP32 __builtin_bswap32
#endif

#include "Ash.h"
#include <iostream>
#include <algorithm>

bool Ash::isAshCompressed(std::vector<uint8_t>& in)
{
	return in.size() > 0x10 &&
		(BYTESWAP32(*reinterpret_cast<uint32_t*>(in.data())) & 0xFFFFFF00) == 0x41534800;

	/* GO VERSION
	return uint32(len(data)) > 0x10 &&
		(BYTESWAP32(*(*uint32)(unsafe.Pointer(&data[0]))) & 0xFFFFFF00) == 0x41534800
	*/
}

std::vector<uint8_t> Ash::decompress(std::vector<uint8_t>& in)
{
	if (!isAshCompressed(in))
	{
		std::cout << "Ash is not compressed" << std::endl;
		return std::vector<uint8_t>();
	}

	uint32_t r[32];
	uint32_t count;

	const auto memAddr = reinterpret_cast<uint64_t>(in.data()); // in
	r[4] = 0x8000000;
	const int64_t inDiff = memAddr - r[4]; // hack to support higher memory addresses than crediar's version

	r[5] = 0x415348;
	r[6] = 0x415348;

	r[5] = BYTESWAP32(*reinterpret_cast<uint32_t*>((r[4] + inDiff + 4)));
	r[5] = r[5] & 0x00FFFFFF;

	const auto size = r[5];
	// std::cout << "Decompressed size: " << size << std::endl;

	std::vector<uint8_t> crap2(size, '\0');
	if (static_cast<uint32_t>(crap2.size()) != size)
	{
		std::cout << "Out of memory 1" << std::endl;
		return std::vector<uint8_t>();
	}
	const auto memAddr2 = reinterpret_cast<uint64_t>(crap2.data()); // outbuf
	r[3] = 0x9000000;
	const int64_t outDiff = memAddr2 - r[3];

	const auto o = r[3];

	r[24] = 0x10;
	r[28] = BYTESWAP32(*reinterpret_cast<uint32_t*>((r[4] + 8 + inDiff)));
	r[25] = 0;
	r[29] = 0;
	r[26] = BYTESWAP32(*reinterpret_cast<uint32_t*>((r[4] + 0xC + inDiff)));
	r[30] = BYTESWAP32(*reinterpret_cast<uint32_t*>((r[4] + r[28] + inDiff)));
	r[28] = r[28] + 4;

	// HACK, pointer to RAM
	std::vector<uint8_t> crap3(std::max(size, 0x100000u), '\0');
	if (crap3.size() != std::max(size, 0x100000u))
	{
		std::cout << "Out of memory 2" << std::endl;
		return std::vector<uint8_t>();
	}

	const auto memAddr3 = reinterpret_cast<uint64_t>(crap3.data()); // outbuf
	r[8] = 0x84000000;
	const int64_t outDiff2 = memAddr3 - r[8];
	memset(reinterpret_cast<void*>(r[8] + outDiff2), 0, 0x100000);
	r[9] = r[8] + 0x07FE;
	r[10] = r[9] + 0x07FE;
	r[11] = r[10] + 0x1FFE;
	r[31] = r[11] + 0x1FFE;
	r[23] = 0x200;
	r[22] = 0x200;
	r[27] = 0;

loc_81332124:

	if (r[25] != 0x1F)
		goto loc_81332140;

	r[0] = r[26] >> 31;
	r[26] = BYTESWAP32(*reinterpret_cast<uint32_t*>((r[4] + r[24] + inDiff)));
	r[25] = 0;
	r[24] = r[24] + 4;
	goto loc_8133214C;

loc_81332140:

	r[0] = r[26] >> 31;
	r[25] = r[25] + 1;
	r[26] = r[26] << 1;

loc_8133214C:

	if (r[0] == 0)
		goto loc_81332174;

	r[0] = r[23] | 0x8000;
	*reinterpret_cast<uint16_t *>(r[31] + outDiff2) = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(r[0])));
	r[0] = r[23] | 0x4000;
	*reinterpret_cast<uint16_t *>(r[31] + 2 + outDiff2) = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(r[0])));

	r[31] = r[31] + 4;
	r[27] = r[27] + 2;
	r[23] = r[23] + 1;
	r[22] = r[22] + 1;

	goto loc_81332124;

loc_81332174:

	r[12] = 9;
	r[21] = r[25] + r[12];
	auto t = r[21];
	if (r[21] > 0x20)
		goto loc_813321AC;

	r[21] = (~(r[12] - 0x20)) + 1;
	r[6] = r[26] >> r[21];
	if (t == 0x20)
		goto loc_8133219C;

	r[26] = r[26] << r[12];
	r[25] = r[25] + r[12];
	goto loc_813321D0;

loc_8133219C:

	r[26] = BYTESWAP32(*reinterpret_cast<uint32_t*>((r[4] + r[24] + inDiff)));
	r[25] = 0;
	r[24] = r[24] + 4;
	goto loc_813321D0;

loc_813321AC:

	r[0] = (~(r[12] - 0x20)) + 1;
	r[6] = r[26] >> r[0];
	r[26] = BYTESWAP32(*reinterpret_cast<uint32_t*>((r[4] + r[24] + inDiff)));
	r[0] = (~(r[21] - 0x40)) + 1;
	r[24] = r[24] + 4;
	r[0] = r[26] >> r[0];
	r[6] = r[6] | r[0];
	r[25] = r[21] - 0x20;
	r[26] = r[26] << r[25];

loc_813321D0:

	r[12] = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(*reinterpret_cast<uint16_t *>((r[31] + outDiff2) - 2))));
	r[31] -= 2;
	r[27] = r[27] - 1;
	r[0] = r[12] & 0x8000;
	r[12] = (r[12] & 0x1FFF) << 1;
	if (r[0] == 0)
		goto loc_813321F8;

	*reinterpret_cast<uint16_t *>(r[9] + r[12] + outDiff2) = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(r[6]))); // ?????
	r[6] = (r[12] & 0x3FFF) >> 1; // extrwi %r6, %r12, 14,17
	if (r[27] != 0)
		goto loc_813321D0;

	goto loc_81332204;

loc_813321F8:

	*reinterpret_cast<uint16_t *>(r[8] + r[12] + outDiff2) = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(r[6])));
	r[23] = r[22];
	goto loc_81332124;

loc_81332204:

	r[23] = 0x800;
	r[22] = 0x800;

loc_8133220C:

	if (r[29] != 0x1F)
		goto loc_81332228;

	r[0] = r[30] >> 31;
	r[30] = BYTESWAP32(*reinterpret_cast<uint32_t*>((r[4] + r[28] + inDiff)));
	r[29] = 0;
	r[28] = r[28] + 4;
	goto loc_81332234;

loc_81332228:

	r[0] = r[30] >> 31;
	r[29] = r[29] + 1;
	r[30] = r[30] << 1;

loc_81332234:

	if (r[0] == 0)
		goto loc_8133225C;

	r[0] = r[23] | 0x8000;
	*reinterpret_cast<uint16_t *>(r[31] + outDiff2) = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(r[0])));
	r[0] = r[23] | 0x4000;
	*reinterpret_cast<uint16_t *>(r[31] + 2 + outDiff2) = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(r[0])));

	r[31] = r[31] + 4;
	r[27] = r[27] + 2;
	r[23] = r[23] + 1;
	r[22] = r[22] + 1;

	goto loc_8133220C;

loc_8133225C:

	r[12] = 0xB;
	r[21] = r[29] + r[12];
	t = r[21];
	if (r[21] > 0x20)
		goto loc_81332294;

	r[21] = (~(r[12] - 0x20)) + 1;
	r[7] = r[30] >> r[21];
	if (t == 0x20)
		goto loc_81332284;

	r[30] = r[30] << r[12];
	r[29] = r[29] + r[12];
	goto loc_813322B8;

loc_81332284:

	r[30] = BYTESWAP32(*reinterpret_cast<uint32_t *>(r[4] + r[28] + inDiff));
	r[29] = 0;
	r[28] = r[28] + 4;
	goto loc_813322B8;

loc_81332294:

	r[0] = (~(r[12] - 0x20)) + 1;
	r[7] = r[30] >> r[0];
	r[30] = BYTESWAP32(*reinterpret_cast<uint32_t *>(r[4] + r[28] + inDiff));
	r[0] = (~(r[21] - 0x40)) + 1;
	r[28] = r[28] + 4;
	r[0] = r[30] >> r[0];
	r[7] = r[7] | r[0];
	r[29] = r[21] - 0x20;
	r[30] = r[30] << r[29];

loc_813322B8:

	r[12] = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(*reinterpret_cast<uint16_t *>((r[31] + outDiff2) - 2))));
	r[31] -= 2;
	r[27] = r[27] - 1;
	r[0] = r[12] & 0x8000;
	r[12] = (r[12] & 0x1FFF) << 1;
	if (r[0] == 0)
		goto loc_813322E0;

	*reinterpret_cast<uint16_t *>(r[11] + r[12] + outDiff2) = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(r[7]))); // ????
	r[7] = (r[12] & 0x3FFF) >> 1; // extrwi %r7, %r12, 14,17
	if (r[27] != 0)
		goto loc_813322B8;

	goto loc_813322EC;

loc_813322E0:

	*reinterpret_cast<uint16_t *>(r[10] + r[12] + outDiff2) = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(r[7])));
	r[23] = r[22];
	goto loc_8133220C;

loc_813322EC:

	r[0] = r[5];

loc_813322F0:

	r[12] = r[6];

loc_813322F4:

	if (r[12] < 0x200)
		goto loc_8133233C;

	if (r[25] != 0x1F)
		goto loc_81332318;

	r[31] = r[26] >> 31;
	r[26] = BYTESWAP32(*reinterpret_cast<uint32_t *>(r[4] + r[24] + inDiff));
	r[24] = r[24] + 4;
	r[25] = 0;
	goto loc_81332324;

loc_81332318:

	r[31] = r[26] >> 31;
	r[25] = r[25] + 1;
	r[26] = r[26] << 1;

loc_81332324:

	r[27] = r[12] << 1;
	if (r[31] != 0)
		goto loc_81332334;

	r[12] = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(*reinterpret_cast<uint16_t *>(r[8] + r[27] + outDiff2))));
	goto loc_813322F4;

loc_81332334:

	r[12] = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(*reinterpret_cast<uint16_t *>(r[9] + r[27] + outDiff2))));
	goto loc_813322F4;

loc_8133233C:

	if (r[12] >= 0x100)
		goto loc_8133235C;

	*reinterpret_cast<uint8_t *>(r[3] + outDiff) = r[12];
	r[3] = r[3] + 1;
	r[5] = r[5] - 1;
	if (r[5] != 0)
		goto loc_813322F0;

	goto loc_81332434;

loc_8133235C:

	r[23] = r[7];

loc_81332360:

	if (r[23] < 0x800)
		goto loc_813323A8;

	if (r[29] != 0x1F)
		goto loc_81332384;

	r[31] = r[30] >> 31;
	r[30] = BYTESWAP32(*reinterpret_cast<uint32_t *>(r[4] + r[28] + inDiff));
	r[28] = r[28] + 4;
	r[29] = 0;
	goto loc_81332390;

loc_81332384:

	r[31] = r[30] >> 31;
	r[29] = r[29] + 1;
	r[30] = r[30] << 1;

loc_81332390:

	r[27] = r[23] << 1;
	if (r[31] != 0)
		goto loc_813323A0;

	r[23] = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(*reinterpret_cast<uint16_t *>(r[10] + r[27] + outDiff2)))
		);
	goto loc_81332360;

loc_813323A0:

	r[23] = static_cast<uint16_t>(BYTESWAP16(static_cast<uint16_t>(*reinterpret_cast<uint16_t *>(r[11] + r[27] + outDiff2)))
		);
	goto loc_81332360;

loc_813323A8:

	r[12] = r[12] - 0xFD;
	r[23] = ~r[23] + r[3] + 1;
	r[5] = ~r[12] + r[5] + 1;
	r[31] = r[12] >> 3;

	if (r[31] == 0)
		goto loc_81332414;

	count = r[31];

loc_813323C0:

	r[31] = *reinterpret_cast<uint8_t *>((r[23] + outDiff) - 1);
	*reinterpret_cast<uint8_t *>(r[3] + outDiff) = r[31];

	r[31] = *reinterpret_cast<uint8_t *>(r[23] + outDiff);
	*reinterpret_cast<uint8_t *>(r[3] + 1 + outDiff) = r[31];

	r[31] = *reinterpret_cast<uint8_t *>((r[23] + outDiff) + 1);
	*reinterpret_cast<uint8_t *>(r[3] + 2 + outDiff) = r[31];

	r[31] = *reinterpret_cast<uint8_t *>((r[23] + outDiff) + 2);
	*reinterpret_cast<uint8_t *>(r[3] + 3 + outDiff) = r[31];

	r[31] = *reinterpret_cast<uint8_t *>((r[23] + outDiff) + 3);
	*reinterpret_cast<uint8_t *>(r[3] + 4 + outDiff) = r[31];

	r[31] = *reinterpret_cast<uint8_t *>((r[23] + outDiff) + 4);
	*reinterpret_cast<uint8_t *>(r[3] + 5 + outDiff) = r[31];

	r[31] = *reinterpret_cast<uint8_t *>((r[23] + outDiff) + 5);
	*reinterpret_cast<uint8_t *>(r[3] + 6 + outDiff) = r[31];

	r[31] = *reinterpret_cast<uint8_t *>((r[23] + outDiff) + 6);
	*reinterpret_cast<uint8_t *>(r[3] + 7 + outDiff) = r[31];

	r[23] = r[23] + 8;
	r[3] = r[3] + 8;

	if (--count)
		goto loc_813323C0;

	r[12] = r[12] & 7;
	if (r[12] == 0)
		goto loc_8133242C;

loc_81332414:

	count = r[12];

loc_81332418:

	r[31] = *reinterpret_cast<uint8_t *>((r[23] + outDiff) - 1);
	r[23] = r[23] + 1;
	*reinterpret_cast<uint8_t *>(r[3] + outDiff) = r[31];
	r[3] = r[3] + 1;

	if (--count)
		goto loc_81332418;

loc_8133242C:

	if (r[5] != 0)
		goto loc_813322F0;

loc_81332434:

	r[3] = r[0];

	std::vector<uint8_t> ret(reinterpret_cast<const char*>(o + outDiff), reinterpret_cast<const char*>(o + outDiff) + r[3]);

	return ret;
}
