package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.CartInfo;
import com.seamlance.perfume.mapper.CartMapper;
import com.seamlance.perfume.service.CartService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class CartDelegateImpl implements CartDelegate{

    @Autowired
    private CartService cartService;

    @Autowired
    private CartMapper cartMapper;

    @Override
    public CartInfo addCart(CartInfo cartInfo) {
        return cartMapper.toInfo(cartService.addCart(cartMapper.toEntity(cartInfo)));
    }

    @Override
    public List<CartInfo> getAll() {
        return cartMapper.toInfoList(cartService.getAll());
    }

    @Override
    public void deleteCart(String cartId) {
        cartService.deleteCart(cartId);
    }

    @Override
    public CartInfo updateCart(String cartId, CartInfo cartInfo) {
        return cartMapper.toInfo(cartService.updateCart(cartId, cartMapper.toEntity(cartInfo)));
    }
}
