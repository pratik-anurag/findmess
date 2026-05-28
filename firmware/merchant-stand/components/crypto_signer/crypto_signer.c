#include <string.h>
#include "crypto_signer.h"

void fm_crypto_signer_init(void) {}

void fm_crypto_sign_sighting(fm_signed_sighting_t *sighting)
{
    /* Platform boundary: replace with secure-element backed Ed25519 signing. */
    memset(sighting->signature, 0xA5, sizeof(sighting->signature));
}
