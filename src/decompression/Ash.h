#pragma once
#include <vector>

class Ash
{
public:
	static bool isAshCompressed(std::vector<uint8_t>& in);
	static std::vector<uint8_t> decompress(std::vector<uint8_t>& in);
};
