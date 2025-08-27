package com.seamlance.perfume.security;

public class JwtConstant {
    public static final long EXPIRATION_TIME_IN_MILLISECONDS = 1000L * 60L * 60L; // expiration time 1 hour
    public static final String SECRET = "HpnchpoAmckdaLGJCLCDOLIVALSDCKJHFKPerfumeklcasdiuoeaUNDSKJCELOVScksdcGmdGLSK";
    public static final String ALGORITHM = "HmacSHA256";
    public static final String TOKEN_PREFIX = "Bearer "; // There must be space after bearer
    public static final String HEADER_STRING = "Authorization";
}
