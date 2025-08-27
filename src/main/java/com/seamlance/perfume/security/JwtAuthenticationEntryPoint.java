package com.seamlance.perfume.security;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.dto.ApiResponse;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
public class JwtAuthenticationEntryPoint implements AuthenticationEntryPoint {
    public static final String ERROR_MASSAGE = "You are not allowed to perform this action";
    private ObjectMapper mapper = new ObjectMapper();
    @Override
    public void commence(HttpServletRequest request, HttpServletResponse response, AuthenticationException authException) throws IOException, ServletException {
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.setContentType(Constant.APPLICATION_JSON);

        ApiResponse apiResponse = new ApiResponse(true, ERROR_MASSAGE, null);
        response.getWriter().write(mapper.writeValueAsString(apiResponse));
    }
}
