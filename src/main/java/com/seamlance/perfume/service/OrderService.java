package com.seamlance.perfume.service;

import com.seamlance.perfume.dto.OrderRequest;
import com.seamlance.perfume.entity.Order;
import com.seamlance.perfume.enums.OrderStatus;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

import java.time.LocalDateTime;
import java.util.List;

public interface OrderService {
    Order placeOrder(OrderRequest orderRequest);
    void updateOrderItemStatus(String orderId, String status);
    Page<Order> filterOrderItem(OrderStatus status, LocalDateTime startDate, LocalDateTime endDate, Pageable pageable);
    List<Order> getOrderHistoryByUser(String userId);
}
