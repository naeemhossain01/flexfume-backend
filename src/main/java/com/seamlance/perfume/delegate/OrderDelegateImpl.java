package com.seamlance.perfume.delegate;

import com.seamlance.perfume.dto.OrderRequest;
import com.seamlance.perfume.dto.OrderResponse;
import com.seamlance.perfume.entity.Order;
import com.seamlance.perfume.enums.OrderStatus;
import com.seamlance.perfume.info.OrderInfo;
import com.seamlance.perfume.mapper.OrderMapper;
import com.seamlance.perfume.service.OrderService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Component;

import java.time.LocalDateTime;
import java.util.List;
import java.util.stream.Collectors;

@Component
public class OrderDelegateImpl implements OrderDelegate {
    @Autowired
    private OrderService orderService;

    @Autowired
    private OrderMapper orderMapper;

    @Override
    public OrderInfo placeOrder(OrderRequest orderRequest) {
        return orderMapper.toInfo(orderService.placeOrder(orderRequest));
    }

    @Override
    public void updateOrderStatus(String orderId, String status) {
        orderService.updateOrderItemStatus(orderId, status);
    }

    @Override
    public OrderResponse filterOrder(OrderStatus status, LocalDateTime startDate, LocalDateTime endDate, Pageable pageable) {
        Page<Order> orderPage = orderService.filterOrderItem(status, startDate, endDate, pageable);

        List<OrderInfo> orderInfoList = orderPage.getContent().stream().map(orderMapper::toInfo).collect(Collectors.toList());

        OrderResponse orderResponse = new OrderResponse();
        orderResponse.setOrderInfoList(orderInfoList);
        orderResponse.setTotalPage(orderPage.getTotalPages());
        orderResponse.setTotalElements(orderPage.getTotalElements());

        return orderResponse;
    }

    @Override
    public List<OrderInfo> getOrderHistoryByUser(String userId) {
        return orderMapper.toInfoList(orderService.getOrderHistoryByUser(userId));
    }
}
