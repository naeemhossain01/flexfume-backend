package com.seamlance.perfume.security;


import java.nio.charset.StandardCharsets;
import java.util.Date;
import java.util.function.Function;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;

import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import com.seamlance.perfume.entity.User;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import jakarta.annotation.PostConstruct;

@Service
public class JwtUtils {
    private SecretKey key;

    @PostConstruct
    private void init() {
        byte[] keyBytes = JwtConstant.SECRET.getBytes(StandardCharsets.UTF_8);
        this.key = new SecretKeySpec(keyBytes, JwtConstant.ALGORITHM);
    }

    public String generateToken(User user){
        String username = user.getPhoneNumber();
        String authorities = user.getRole();
        return generateToken(username, authorities);
    }

    public String generateToken(String username){
        return generateToken(username, "USER");
    }

    public String generateToken(String username, String authorities){
        return Jwts.builder()
                .subject(username)
                .claim("authorities", authorities)
                .issuedAt(new Date(System.currentTimeMillis()))
                .expiration(new Date(System.currentTimeMillis() + JwtConstant.EXPIRATION_TIME_IN_MILLISECONDS))
                .signWith(key)
                .compact();
    }

    public String getUsernameFromToken(String token){
        return extractClaims(token, Claims::getSubject);
    }

    public String getAuthoritiesFromToken(String token){
        return extractClaims(token, claims -> claims.get("authorities", String.class));
    }

    private <T> T extractClaims(String token, Function<Claims, T> claimsTFunction){
        return claimsTFunction.apply(Jwts.parser().verifyWith(key).build().parseSignedClaims(token).getPayload());
    }

    public boolean isTokenValid(String token, UserDetails userDetails){
        final String username = getUsernameFromToken(token);
        return (username.equals(userDetails.getUsername()) && !isTokenExpired(token));
    }

    private boolean isTokenExpired(String token){
        return extractClaims(token, Claims::getExpiration).before(new Date());
    }
}
