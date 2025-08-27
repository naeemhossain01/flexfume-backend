package com.seamlance.perfume.delegate;

import com.seamlance.perfume.dto.OrderRequest;
import com.seamlance.perfume.dto.OrderResponse;
import com.seamlance.perfume.enums.OrderStatus;
import com.seamlance.perfume.info.OrderInfo;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

import java.time.LocalDateTime;
import java.util.List;

public interface OrderDelegate {
    OrderInfo placeOrder(OrderRequest orderRequest);
    void updateOrderStatus(String orderId, String status);
    OrderResponse filterOrder(OrderStatus status, LocalDateTime startDate, LocalDateTime endDate, Pageable pageable);
    List<OrderInfo> getOrderHistoryByUser(String userId);
}
