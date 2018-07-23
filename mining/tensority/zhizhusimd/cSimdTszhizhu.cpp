#include <iostream>
#include <cstdio>
#include <map>
#include <mutex>
#include "cSimdTs.h"
#include "BytomPoW.h"
#include "seed.h"

using namespace std;

uint8_t *SimdTs2impSimd(uint8_t blockheader[32], uint8_t seed[32],uint8_t* resultRes,bool useLock);

uint8_t resultarrSimd[20000][32] = {0};

uint8_t result_lockSimd[32] = {0};
mutex result_lockmutexSimd;

map <vector<uint8_t>, BytomMatList16_Simd*> seedCacheSimd;
mutex seedCacheSimdmutex;

static const int cacheSizeSimd = 10; //"Answer to the Ultimate Question of Life, the Universe, and Everything"
static const int maxidcountSimd = 1000000;
static uint32_t idCountSimd = 0;

uint8_t *SimdTsNoLockSimd(uint8_t blockheader[32], uint8_t seed[32],uint32_t mode){
    //std::cout<<"poollog SimdTsNoLockSimd mode="<< mode  <<std::endl;
    return SimdTs2impSimd(blockheader,seed,resultarrSimd[mode],false);
}

uint8_t *SimdTsWithLockSimd(uint8_t blockheader[32], uint8_t seed[32],uint32_t mode){
    //std::cout<<"poollog SimdTsWithLockSimd mode="<< mode  <<std::endl;
    return SimdTs2impSimd(blockheader,seed,result_lockSimd,true);
}

uint8_t *SimdTs2impSimd(uint8_t blockheader[32], uint8_t seed[32],uint8_t* resultRes,bool useLock){
    vector<uint8_t> seedVec(seed, seed + 32);
    BytomMatList16_Simd* matList_int16_simd;

    seedCacheSimdmutex.lock();
    if(seedCacheSimd.find(seedVec) != seedCacheSimd.end()) {
        // printf("\t---%s---\n", "Seed already exists in the cache.");
        matList_int16_simd = seedCacheSimd[seedVec];
    } else {
        uint32_t exted[32];
        extend(exted, seed); // extends seed to exted
        Words32 extSeed;
        init_seed(extSeed, exted);

        matList_int16_simd = new BytomMatList16_Simd;
        matList_int16_simd->init(extSeed);

        idCountSimd++;
        matList_int16_simd->setId(idCountSimd);

        seedCacheSimd.insert(pair<vector<uint8_t>, BytomMatList16_Simd*>(seedVec, matList_int16_simd));

        if(seedCacheSimd.size() > cacheSizeSimd) {
            for(map<vector<uint8_t>, BytomMatList16_Simd*>::iterator it=seedCacheSimd.begin(); it!=seedCacheSimd.end(); ++it){
                uint32_t id_ = it->second->getId();

                uint32_t step = 0;
                if(idCountSimd < id_ && idCountSimd+maxidcountSimd >= id_ + cacheSizeSimd){
                    delete it->second;
                    seedCacheSimd.erase(it);

                    //std::cout<<"SimdTs2imp.remove  idCount < id_ && idCount+maxidcount >= id_ + cacheSize2 id_="<< id_ <<"  idCount=" << idCount <<std::endl;
                }else if (idCountSimd >= id_+ cacheSizeSimd){
                    delete it->second;
                    seedCacheSimd.erase(it);

                    //std::cout<<"SimdTs2imp.remove  idCount >= id_+ cacheSize2 id_="<< id_ <<"  idCount=" << idCount <<std::endl;
                }
            }

            if(idCountSimd>=maxidcountSimd){
                idCountSimd = 0;
                //std::cout<<"SimdTs2imp.clean idcount  maxidcount="<< maxidcount <<"  idCount=" << idCount <<std::endl;
            }
        }
    }

    seedCacheSimdmutex.unlock();

    if (useLock) {
        result_lockmutexSimd.lock();
    }

    iter_mineBytom(blockheader, 32, resultRes,matList_int16_simd);

    if (useLock) {
        result_lockmutexSimd.unlock();
    }

    return resultRes;
}
