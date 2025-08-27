package com.seamlance.perfume.security;

import java.io.IOException;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.filter.OncePerRequestFilter;

import com.seamlance.perfume.service.CustomUserDetailsService;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

@Component
public class JwtAuthenticationFilter extends OncePerRequestFilter {

    @Autowired
    private JwtUtils jwtUtils;

    @Autowired
    private CustomUserDetailsService customUserDetailsService;

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain) throws ServletException, IOException {
        // Skip JWT processing for public endpoints
        String requestPath = request.getRequestURI();
        if (isPublicEndpoint(requestPath)) {
            filterChain.doFilter(request, response);
            return;
        }
        
        String token = getTokenFromRequest(request);

        if (token != null){
            String username = jwtUtils.getUsernameFromToken(token);

            UserDetails userDetails = customUserDetailsService.loadUserByUsername(username);

            if (StringUtils.hasText(username) && jwtUtils.isTokenValid(token, userDetails)) {

                UsernamePasswordAuthenticationToken authenticationToken = new UsernamePasswordAuthenticationToken(
                        userDetails, null, userDetails.getAuthorities()
                );
                authenticationToken.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));

                SecurityContextHolder.getContext().setAuthentication(authenticationToken);
            } else {
                System.out.println("Invalid header value");
            }

        }
        filterChain.doFilter(request, response);
    }
    
    private boolean isPublicEndpoint(String requestPath) {
        return requestPath.startsWith("/api/v1/auth/") ||
               requestPath.startsWith("/api/v1/product/") ||
               requestPath.startsWith("/api/v1/category/") ||
               requestPath.startsWith("/api/v1/coupon/") ||
               requestPath.startsWith("/api/v1/delivery-cost/") ||
               requestPath.startsWith("/api/v1/guest-cart/") ||
               requestPath.equals("/api/v1/guest-cart") ||
               requestPath.startsWith("/api/v1/checkout/") ||
               requestPath.startsWith("/api/v1/user/reset-password");
    }

    private String getTokenFromRequest(HttpServletRequest request){
        String token = request.getHeader(JwtConstant.HEADER_STRING);
        if (StringUtils.hasText(token) && StringUtils.startsWithIgnoreCase(token, JwtConstant.TOKEN_PREFIX)){
            return token.substring(7);
        }
        return null;
    }
}
