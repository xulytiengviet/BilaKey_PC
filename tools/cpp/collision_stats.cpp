#include <fstream>
#include <iostream>
#include <string>

int main(int argc, char** argv) {
    if (argc != 2) {
        std::cerr << "usage: collision_stats <audit.tsv>\n";
        return 2;
    }
    std::ifstream in(argv[1]);
    if (!in) {
        std::cerr << "cannot open: " << argv[1] << "\n";
        return 3;
    }
    std::string line;
    std::getline(in, line); // header
    std::size_t rows = 0, critical = 0, canonical = 0, explicit_policy = 0;
    while (std::getline(in, line)) {
        if (line.empty()) continue;
        ++rows;
        if (line.find("CRITICAL_PATCH_COLLISION") != std::string::npos) ++critical;
        if (line.find("CANONICAL_AMBIGUITY") != std::string::npos) ++canonical;
        if (line.find("\tEXPLICIT") != std::string::npos) ++explicit_policy;
    }
    std::cout << "rows=" << rows << " critical=" << critical
              << " canonical=" << canonical
              << " explicit_policy=" << explicit_policy << "\n";
    return rows == 56 && critical == 5 && canonical == 51 && explicit_policy == 56 ? 0 : 4;
}
