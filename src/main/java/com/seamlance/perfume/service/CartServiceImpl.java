package com.seamlance.perfume.service;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.stereotype.Service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.CartDao;
import com.seamlance.perfume.entity.Cart;
import com.seamlance.perfume.entity.User;
import com.seamlance.perfume.exception.InvalidRequestsException;
import com.seamlance.perfume.exception.ResourceNotFoundException;

@Service
public class CartServiceImpl implements CartService {
    @Autowired
    private CartDao cartDao;

    @Autowired
    private UserService userService;

    @Override
    public Cart addCart(Cart cart) {
        User user = userService.getLoginUser();
        // Check if the product already exists in the user's cart
        Cart existingCart = cartDao.findByUserIdAndProductId(user.getId(), cart.getProduct().getId());

        if (existingCart != null) {
            // Product already exists, update the quantity
            existingCart.setQuantity(existingCart.getQuantity() + cart.getQuantity());
            try {
                existingCart = cartDao.saveAndFlush(existingCart);
            } catch (DataAccessException e) {
                e.printStackTrace();
            }
            return existingCart;
        } else {
            // Product doesn't exist, create new cart entry
            try {
                cart.setUser(user);
                cart = cartDao.saveAndFlush(cart);
            } catch (DataAccessException e) {
                e.printStackTrace();
            }
            return cart;
        }
    }

    @Override
    public List<Cart> getAll() {
        User user = userService.getLoginUser();

        List<Cart> carts = new ArrayList<>();
        try {
            carts = cartDao.findByUserId(user.getId());
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }

        return carts;
    }

    @Override
    public void deleteCart(String cartId) {
        Cart cart = cartDao.findById(cartId).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.CART_NOT_FOUND));

        User user = userService.getLoginUser();

        if(!cart.getUser().getId().equals(user.getId())) {
            throw new InvalidRequestsException(ErrorConstant.NOT_ALLOWED);
        }

        try {
            cartDao.deleteById(cartId);
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }

    @Override
    public Cart updateCart(String cartId, Cart cart) {
        Cart updatedCart = cartDao.findById(cartId).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.CART_NOT_FOUND));

        int quantity = cart.getQuantity();

        if(quantity <= 0) {
            // Delete the cart item if quantity is 0 or negative
            this.deleteCart(cartId);
            return updatedCart; // Return the original cart for reference
        }

        try {
            updatedCart.setQuantity(quantity);
            updatedCart = cartDao.save(updatedCart);
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }

        return updatedCart;
    }

    @Override
    public void deleteCartByProductId(String productId) {
        User user = userService.getLoginUser();

        Cart cart = cartDao.findByUserIdAndProductId(user.getId(), productId);

        if(cart == null) {
            throw new InvalidRequestsException(ErrorConstant.CART_NOT_FOUND);
        }

        try {
            cartDao.deleteById(cart.getId());
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }

    @Override
    public List<Object> getCartAndProduct(List<String> cartIds) {
        List<Object> carts = new ArrayList<>();

        try {
            carts = cartDao.findByCartIdIn(cartIds);
        } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }

        return carts;
    }
}
