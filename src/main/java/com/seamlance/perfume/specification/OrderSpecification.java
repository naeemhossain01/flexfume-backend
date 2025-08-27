package com.seamlance.perfume.specification;

import com.seamlance.perfume.entity.Order;
import com.seamlance.perfume.enums.OrderStatus;
import org.springframework.data.jpa.domain.Specification;

import java.time.LocalDateTime;

public class OrderSpecification {

    /**
     * To filter order by status
     * @param status
     * @return Order
     */
    public static Specification<Order> hasStatus(OrderStatus status) {
        return ((root, query, criteriaBuilder) ->
            status != null ? criteriaBuilder.equal(root.get("status"), status) : null);
    }


    /**
     * TO filter order by date
     * @param startDate
     * @param endDate
     * @return Order
     */
    public static Specification<Order> createBetween(LocalDateTime startDate, LocalDateTime endDate) {
        return  ((root, query, criteriaBuilder) ->{
            if(startDate != null && endDate != null) {
                return criteriaBuilder.between(root.get("createdDate"), startDate, endDate);
            } else if(startDate != null) {
                return criteriaBuilder.greaterThanOrEqualTo(root.get("createdDate"), startDate);
            } else if(endDate != null) {
                return criteriaBuilder.lessThanOrEqualTo(root.get("createdDate"), endDate);
            }

            return  null;
        });
    }
}
