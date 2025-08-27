package com.seamlance.perfume.dto;

public record ApiResponse<T> (Boolean error, String message, T response) {
    public static final String MESSAGE_SUCCESS = "SUCCESS";

    public static <T> ApiResponse<T> success(T response) {
        return new ApiResponse<>(false, MESSAGE_SUCCESS, response);
    }

    public static <T> ApiResponse<T> error(String message) {
        return new ApiResponse<>(true, message, null);
    }
}
