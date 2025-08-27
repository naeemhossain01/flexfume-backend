package com.seamlance.perfume.dao;

import com.seamlance.perfume.entity.Product;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface ProductDao extends JpaRepository<Product, String> {
    List<Product> findByCategoryId(String categoryId);
    List<Product> findByProductNameContainingOrDescriptionContaining(String name, String description);
}
