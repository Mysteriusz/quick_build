#if !defined(AX_ERROR_CODE_INT)
#define AX_ERROR_CODE_INT

/*
 
   	Code mapping:

   	axres bit 0-11 		: Error code (MAX -> 0xfff)
   	axres bit 12-15		: Reserved
   	axres bit 16-31		: Meta error bits
   	axres bit 32-63		: Meta error

	Further bit fields may be added
   	
	Code structure can be revealed by casting 
		- (axres_s*)&axres

*/

#include "ax_type.h"

#define AX_META_NTSTATUS	(astp(axres, (axres_s){.meta.ntstatus = true}))

// ERROR CODES CANNOT BE BIGGER THAN 0xFFF (12 bits) 
#define AX_SUCC 		(astp(axres, ((axres_s){.err = 0x00, .meta = {0}})))

// "INVALID" codes

#define AX_INV_ARG 		(astp(axres, ((axres_s){.err = 0x01, .meta = {0}})))
#define AX_INV_ARG_MSG 		u"Invalid argument passed."

#define AX_INV_DATA 		(astp(axres, ((axres_s){.err = 0x02, .meta = {0}})))
#define AX_INV_DATA_MSG 	u"Invalid data passed."

#define AX_INV_BUF 		(astp(axres, ((axres_s){.err = 0x03, .meta = {0}})))
#define AX_INV_BUF_MSG 		u"Invalid buffer passed."

#define AX_INV_CODE 		(astp(axres, ((axres_s){.err = 0x04, .meta = {0}})))
#define AX_INV_CODE_MSG 	u"Invalid code received."

#define AX_INV_FILE 		(astp(axres, ((axres_s){.err = 0x05, .meta = {0}}))) // EXCLUSIVE TO _io_file STRUCTURE ERRORS
#define AX_INV_FILE_MSG 	u"Invalid file structure."

#define AX_INV_ENC 		(astp(axres, ((axres_s){.err = 0x06, .meta = {0}}))) // EXCLUSIVE TO _io_file_enc TYPE ERRORS
#define AX_INV_ENC_MSG 		u"Invalid file encoding."

#define AX_INV_FMT 		(astp(axres, ((axres_s){.err = 0x07, .meta = {0}})))
#define AX_INV_FMT_MSG 		u"Invalid value format."

#define AX_INV_IND 		(astp(axres, ((axres_s){.err = 0x08, .meta = {0}})))
#define AX_INV_IND_MSG 		u"Index out of bounds."

#define AX_INV_MEM 		(astp(axres, ((axres_s){.err = 0x09, .meta = {0}})))
#define AX_INV_MEM_MSG 		u"Memory corruption."

#define AX_INV_PATH 		(astp(axres, ((axres_s){.err = 0x0a, .meta = {0}})))
#define AX_INV_PATH_MSG 	u"Invalid path."

// "BUFFER" codes

#define AX_BUF_TOO_SMALL 	(astp(axres, ((axres_s){.err = 0x10, .meta = {0}})))
#define AX_BUF_TOO_SMALL_MSG 	u"Buffer too small."

#define AX_BUF_TOO_BIG 		(astp(axres, ((axres_s){.err = 0x11, .meta = {0}})))
#define AX_BUF_TOO_BIG_MSG 	u"Buffer too big."

// "NOT" codes

#define AX_NOT_FND 		(astp(axres, ((axres_s){.err = 0x20, .meta = {0}})))
#define AX_NOT_IMP 		(astp(axres, ((axres_s){.err = 0x21, .meta = {0}})))

// "ACCESS" codes

#define AX_ACC_DEN		(astp(axres, ((axres_s){.err = 0x40, .meta = {0}})))

// "UNKNOWN" codes

#define AX_UNK_ERR 		(astp(axres, ((axres_s){.err = 0x50, .meta = {0}})))

// "USER" codes

#define AX_USR_ERR		(astp(axres, ((axres_s){.err = 0x60, .meta = {0}})))

// "IR" codes (Exclusive to Intermediate Representation of the MTE)

#define AX_IR_BLOCKED 		(astp(axres, ((axres_s){.err = 0x70, .meta = {0}})))

// "MTE" codes (Exclusive to Micro Translation Engine)

#define AX_MTE_INV_ARCH 	(astp(axres, ((axres_s){.err = 0x80, .meta = {0}})))
#define AX_MTE_INV_INSTR 	(astp(axres, ((axres_s){.err = 0x81, .meta = {0}})))

static inline axres _ax_buf_err(
	u64 		size,
	u64 		buf_size	
){
#if defined(AX_STRICT_BUF_SIZE)
	if (size < buf_size){
		return AX_BUF_TOO_BIG;
	}
#endif
	if (size > buf_size){
		return AX_BUF_TOO_SMALL;
	}

	return AX_SUCC;
}

#endif // !defined(AX_ERROR_CODE_INT)

