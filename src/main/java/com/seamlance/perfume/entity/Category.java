package com.seamlance.perfume.entity;

import com.fasterxml.jackson.annotation.JsonManagedReference;
import com.seamlance.perfume.audit.AbstractAudit;
import jakarta.persistence.*;

import java.util.List;

@Entity
@Table(name = "CATEGORY")
@AttributeOverride(name = "id", column = @Column(name = "CATEGORY_ID"))
public class Category extends AbstractAudit {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "NAME", unique = true)
    private String name;

    @OneToMany(mappedBy = "category", fetch = FetchType.LAZY, cascade = CascadeType.ALL)
    @JsonManagedReference
    private List<Product> productList;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public List<Product> getProductList() {
        return productList;
    }

    public void setProductList(List<Product> productList) {
        this.productList = productList;
    }
}
