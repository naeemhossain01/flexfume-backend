package com.seamlance.perfume.service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.OrderDao;
import com.seamlance.perfume.dto.OrderRequest;
import com.seamlance.perfume.entity.*;
import com.seamlance.perfume.enums.OrderStatus;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import com.seamlance.perfume.specification.OrderSpecification;
import jakarta.transaction.Transactional;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

@Service
public class OrderServiceImpl implements OrderService {

    @Autowired
    private OrderDao orderDao;

    @Autowired
    private UserService userService;

    @Autowired
    private ProductService productService;

    @Autowired
    private CartService cartService;

    @Autowired
    private DeliveryCostService deliveryCostService;

    @Autowired
    private CouponUsageService couponUsageService;

    private static final BigDecimal ONE_HUNDRED = new BigDecimal("100");

    @Override
    @Transactional(rollbackOn = Exception.class)
    public Order placeOrder(OrderRequest orderRequest) {
        User user = userService.getLoginUser();

        List<OrderItem> orderItemList = orderRequest.getOrderItemRequestList().stream().map(orderItemRequest -> {
            Product product =  productService.getProductById(orderItemRequest.getProductId());

            OrderItem orderItem = new OrderItem();
            orderItem.setProduct(product);
            orderItem.setQuantity(orderItemRequest.getQuantity());
            orderItem.setPrice(calculateProductPrice(orderItemRequest.getQuantity(), product.getPrice(), product.getDiscount()));
            return orderItem;
        }).collect(Collectors.toList());

        BigDecimal totalPrice = this.calculateProductTotalCost(orderRequest, orderItemList).add(this.getDeliveryCost(user));

        if(orderRequest.getCouponCode() != null) {
            CouponUsage couponUsage = couponUsageService.getCouponUsage(orderRequest.getCouponCode(), user.getId());
            totalPrice = totalPrice.subtract(couponUsage != null ? couponUsage.getDiscountedAmount() : BigDecimal.ZERO);
        }

        Order order = new Order();
        order.setOrderItemList(orderItemList);
        order.setTotalPrice(totalPrice);
        order.setUser(user);
        order.setStatus(OrderStatus.PENDING);

        Order finalOrder = order;
        orderItemList.forEach(orderItem -> {
            cartService.deleteCartByProductId(orderItem.getProduct().getId());
            orderItem.setOrder(finalOrder);
        });

        try {
            order = orderDao.save(order);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return order;
    }

    @Override
    public void updateOrderItemStatus(String orderId, String status) {
        Order order = orderDao.findById(orderId).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.ORDER_NOT_FOUND));

        order.setStatus(OrderStatus.valueOf(status.toUpperCase()));

        try {
            orderDao.save(order);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }
    }

    @Override
    public Page<Order> filterOrderItem(OrderStatus status, LocalDateTime startDate, LocalDateTime endDate, Pageable pageable) {
        Specification<Order> spec = Specification.where(OrderSpecification.hasStatus(status))
                .and(OrderSpecification.createBetween(startDate, endDate));

        Page<Order> orderPage = orderDao.findAll(spec, pageable);

        if(orderPage.isEmpty()) {
            throw new ResourceNotFoundException(ErrorConstant.ORDER_NOT_FOUND);
        }

        return orderPage;
    }

    @Override
    public List<Order> getOrderHistoryByUser(String userId) {
        User user = userService.getUser(userId);

        if(user == null) {
            throw new ResourceNotFoundException(ErrorConstant.USER_NOT_FOUND);
        }

        List<Order> orders = new ArrayList<>();
        try {
            orders = orderDao.findByUserId(userId);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return orders;
    }

    private BigDecimal calculateProductPrice(int quantity, BigDecimal price, Discount discount) {
        BigDecimal productIndividualPrice = (discount != null && discount.getPercentage() > 0)
                ? price.subtract(calculateDiscountAmount(price, discount.getPercentage()))
                : price;
        return productIndividualPrice.multiply(BigDecimal.valueOf(quantity));
    }

    private BigDecimal calculateProductTotalCost(OrderRequest orderRequest, List<OrderItem> orderItemList) {
        BigDecimal productCost = orderRequest.getTotalPrice() != null && orderRequest.getTotalPrice().compareTo(BigDecimal.ZERO) > 0
                ? orderRequest.getTotalPrice()
                : orderItemList.stream().map(OrderItem::getPrice).reduce(BigDecimal.ZERO, BigDecimal::add);

        return productCost;
    }

    private BigDecimal getDeliveryCost(User user) {
        List<DeliveryCost> deliveryCost = new ArrayList<>();
        if(user != null && user.getAddress() !=null && user.getAddress().getCity() != null) {
           deliveryCost = deliveryCostService.getDeliveryCostByLocation(user.getAddress().getCity());
        }

        if (deliveryCost.isEmpty()) {
            throw  new ResourceNotFoundException(ErrorConstant.INVALID_DELIVERY_LOCATION);
        }

        // Double check for accurate delivery cost
        return deliveryCost.stream()
                .filter(dc -> user.getAddress().getCity().equalsIgnoreCase(dc.getLocation()))
                .map(DeliveryCost::getCost)
                .findFirst()
                .orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.INVALID_DELIVERY_LOCATION));
    }

    /**
     * This function calculate discount amount in round figure
     * @param price
     * @param percentage
     * @return (5.5 => 5, 2.5 => 5, 1.6 => 2, 1.1 => 1)
     */
    private BigDecimal calculateDiscountAmount(BigDecimal price, int percentage) {
        return price.multiply(BigDecimal.valueOf(percentage)).divide(ONE_HUNDRED, 2, RoundingMode.HALF_DOWN);
    }
}
