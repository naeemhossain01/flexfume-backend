package com.seamlance.perfume.controller;

import com.seamlance.perfume.delegate.OrderDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.dto.OrderRequest;
import com.seamlance.perfume.dto.OrderResponse;
import com.seamlance.perfume.enums.OrderStatus;
import com.seamlance.perfume.info.OrderInfo;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.List;

@RestController
@RequestMapping("/api/v1/order")
public class OrderController {

    public static final String ORDER_STATUS = "Order status changes successfully";
    @Autowired
    private OrderDelegate orderDelegate;


    @PostMapping()
    public ApiResponse<OrderInfo> placeOrder(@RequestBody OrderRequest orderRequest) {
        return ApiResponse.success(orderDelegate.placeOrder(orderRequest));
    }

    @PutMapping("/{orderId}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<String> updateOrderStatus(@PathVariable String orderId, @RequestParam String status) {
        orderDelegate.updateOrderStatus(orderId, status);
        return ApiResponse.success(ORDER_STATUS);
    }

    @GetMapping("/filter")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<OrderResponse> filterOrderItems(
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)LocalDateTime startDate,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)LocalDateTime endDate,
            @RequestParam(required = false) String status,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "1000") int size
            ) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.DESC, "id"));
        OrderStatus orderStatus = status != null ? OrderStatus.valueOf(status.toUpperCase()) : null;

        OrderResponse orderResponse = orderDelegate.filterOrder(orderStatus, startDate, endDate, pageable);

        return ApiResponse.success(orderResponse);
    }

    @GetMapping("/history/{userId}")
    public ApiResponse<List<OrderInfo>> getOrderHistory(@PathVariable String userId) {
        return ApiResponse.success(orderDelegate.getOrderHistoryByUser(userId));
    }
}
