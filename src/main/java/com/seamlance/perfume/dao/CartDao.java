package com.seamlance.perfume.dao;

import com.seamlance.perfume.entity.Cart;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface CartDao extends JpaRepository<Cart, String> {
    List<Cart> findByUserId(String userId);
    Cart findByUserIdAndProductId(String userId, String productId);

    @Query(nativeQuery = true, value = "SELECT C.CART_ID, C.QUANTITY, P.PRODUCT_ID, P.PRICE, D.PERCENTAGE AS DISCOUNT" +
            " FROM ((CART C JOIN PRODUCT P ON P.PRODUCT_ID = C.PRODUCT_ID)" +
            " LEFT JOIN DISCOUNT D ON D.PRODUCT_ID = P.PRODUCT_ID) WHERE CART_ID IN (?)")
    List<Object> findByCartIdIn(List<String> cartIds);
}
