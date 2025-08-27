package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.CartInfo;

import java.util.List;

public interface CartDelegate {
    CartInfo addCart(CartInfo cartInfo);
    List<CartInfo> getAll();
    void deleteCart(String cartId);
    CartInfo updateCart(String cartId, CartInfo cartInfo);
}
