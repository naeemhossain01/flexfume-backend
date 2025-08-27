package com.seamlance.perfume.service;

import com.seamlance.perfume.entity.Cart;

import java.util.List;

public interface CartService {
    Cart addCart(Cart cart);
    List<Cart> getAll();
    void deleteCart(String cartId);
    Cart updateCart(String cartId, Cart cart);
    void deleteCartByProductId(String productId);
    List<Object> getCartAndProduct(List<String> cartIds);
}
