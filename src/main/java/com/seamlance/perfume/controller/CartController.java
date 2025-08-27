package com.seamlance.perfume.controller;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.seamlance.perfume.delegate.CartDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.info.CartInfo;

@RestController
@RequestMapping("/api/v1/cart")
public class CartController {

    @Autowired
    private CartDelegate cartDelegate;


    @PostMapping()
    @PreAuthorize("hasAuthority('USER') or hasAuthority('ADMIN')")
    public ApiResponse<CartInfo> addCart(@RequestBody CartInfo cartInfo) {
        return ApiResponse.success(cartDelegate.addCart(cartInfo));
    }

    @GetMapping()
    @PreAuthorize("hasAuthority('USER') or hasAuthority('ADMIN')")
    public ApiResponse<List<CartInfo>> getAll() {
        return ApiResponse.success(cartDelegate.getAll());
    }

    @DeleteMapping("/{id}")
    @PreAuthorize("hasAuthority('USER') or hasAuthority('ADMIN')")
    public ApiResponse<String> deleteCart(@PathVariable String id) {
        cartDelegate.deleteCart(id);
        return ApiResponse.success("CART REMOVED");
    }

    @PutMapping("/{id}")
    @PreAuthorize("hasAuthority('USER') or hasAuthority('ADMIN')")
    public ApiResponse<CartInfo> updateCart(@PathVariable String id, @RequestBody CartInfo cartInfo) {
        return ApiResponse.success(cartDelegate.updateCart(id, cartInfo));
    }
}
