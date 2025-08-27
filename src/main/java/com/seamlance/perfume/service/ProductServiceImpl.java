package com.seamlance.perfume.service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.ProductDao;
import com.seamlance.perfume.entity.Category;
import com.seamlance.perfume.entity.Product;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import jakarta.persistence.PersistenceException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.util.ArrayList;
import java.util.List;

@Service
public class ProductServiceImpl implements ProductService {
    @Autowired
    private ProductDao productDao;

    @Autowired
    private CategoryService categoryService;

    @Autowired
    private AwsS3Service awsS3Service;

    @Override
    public Product addProduct(Product product) {
        if(product == null) {
            // throw exception
            return  null;
        }
        // TODO: add validation

        Category category = categoryService.getCategoryById(product.getCategory().getId());

        if(category == null) {
            throw new ResourceNotFoundException(ErrorConstant.CATEGORY_NOT_FOUND);
        }


        try {
            product = productDao.saveAndFlush(product);
        } catch (DataAccessException e) {

        }catch (PersistenceException e) {

        }
        return product;
    }

    @Override
    public Product updateProduct(String productId, Product product) {
        if(productId.isEmpty() || product == null) {
            // throw exception here
            return null;
        }

        // TODO: add validation

        Product updatedProduct = productDao.findById(productId).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.PRODUCT_NOT_FOUND));

        if(product.getProductName() != null) updatedProduct.setProductName(product.getProductName());
        if(product.getProductCode() != null) updatedProduct.setProductCode(product.getProductCode());
        if(product.getDescription() != null) updatedProduct.setDescription(product.getDescription());
        if(product.getImageUrl() != null) updatedProduct.setImageUrl(product.getImageUrl());
        if(product.getCategory() != null) updatedProduct.setCategory(product.getCategory());
        if(product.getPrice() != null) updatedProduct.setPrice(product.getPrice());

        if(product.getCategory() != null) {
            Category category = categoryService.getCategoryById(product.getCategory().getId());

            if(category == null) {
                throw new ResourceNotFoundException(ErrorConstant.CATEGORY_NOT_FOUND);
            }
        }



        try {
            updatedProduct = productDao.save(updatedProduct);
        } catch (DataAccessException e) {

        }catch (PersistenceException e) {

        }

        return updatedProduct;
    }

    @Override
    public Product getProductById(String productId) {
        if(productId.isEmpty()) {
            return  null;
        }

        Product product = null;

        product = productDao.findById(productId).orElseThrow();

        return product;
    }

    @Override
    public List<Product> getAllProduct() {
        List<Product> products = new ArrayList<>();

        try {
            products = productDao.findAll();
        } catch (DataAccessException e) {

        } catch (PersistenceException e) {

        }
        return products;
    }

    @Override
    public void deleteProduct(String productId) {
        if(productId.isEmpty()) {
            return;
        }

        productDao.findById(productId)
                .orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.PRODUCT_NOT_FOUND));

        try {
            productDao.deleteById(productId);
        } catch (DataAccessException e) {

        } catch (PersistenceException e) {

        }
    }

    @Override
    public List<Product> getProductByCategory(String categoryId) {
        List<Product> products = new ArrayList<>();

        try {
            products = productDao.findByCategoryId(categoryId);
            if(products.isEmpty()) {
                throw new ResourceNotFoundException(ErrorConstant.NO_PRODUCT_FOUND_BY_CATEGORY);
            }

        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return products;

    }

    @Override
    public List<Product> searchProduct(String value) {
        List<Product> products = new ArrayList<>();

        try {
            products = productDao.findByProductNameContainingOrDescriptionContaining(value, value);
            if(products.isEmpty()) {
                throw new ResourceNotFoundException(ErrorConstant.NO_PRODUCT_FOUND_BY_CATEGORY);
            }

        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return products;
    }

    @Override
    public Product uploadImage(String productId, MultipartFile file) {
        Product product = productDao.findById(productId).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.PRODUCT_NOT_FOUND));

        if(file == null || file.isEmpty()) {
            throw new ResourceNotFoundException(ErrorConstant.IMAGE_NOT_FOUND);
        }

        String imageUrl = awsS3Service.saveImageToS3(file);

        try {
            product.setImageUrl(imageUrl);
            product = productDao.save(product);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return product;
    }
}
