package com.seamlance.perfume.exception;

import com.seamlance.perfume.dto.ApiResponse;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(ResourceNotFoundException.class)
    @ResponseStatus(value = HttpStatus.NOT_FOUND)
    public ApiResponse<String> handleResourceNotFoundException(ResourceNotFoundException ex) {
        return ApiResponse.error(ex.getMessage());
    }

    @ExceptionHandler(InvalidCredentialsException.class)
    @ResponseStatus(value = HttpStatus.BAD_REQUEST)
    public ApiResponse<String> handleInvalidCredentialsException(InvalidCredentialsException ex) {
        return ApiResponse.error(ex.getMessage());
    }

    @ExceptionHandler(InvalidRequestsException.class)
    @ResponseStatus(value = HttpStatus.BAD_REQUEST)
    public ApiResponse<String> handleBadRequestException(InvalidRequestsException ex) {
        return ApiResponse.error(ex.getMessage());
    }

    @ExceptionHandler(InvalidOtpSenderTypeException.class)
    @ResponseStatus(value = HttpStatus.NOT_FOUND)
    public ApiResponse<String> handleInvalidOtpSenderException(InvalidOtpSenderTypeException ex) {
        return ApiResponse.error(ex.getMessage());
    }
}
