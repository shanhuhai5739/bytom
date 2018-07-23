#ifndef _C_SIMD_TENSOR_ZHIZHU_H_
#define _C_SIMD_TENSOR_ZHIZHU_H_

#ifdef __cplusplus
extern "C" {
#endif
	#include <stdint.h>

    uint8_t *SimdTsNoLockSimd(uint8_t blockheader[32], uint8_t seed[32],uint32_t mode);

    uint8_t *SimdTsWithLockSimd(uint8_t blockheader[32], uint8_t seed[32],uint32_t mode);

#ifdef __cplusplus
}
#endif

#endif
