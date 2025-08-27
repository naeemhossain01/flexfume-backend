package com.seamlance.perfume.entity;

import com.seamlance.perfume.audit.AbstractAudit;
import jakarta.persistence.*;

@Entity
@Table(name = "DISCOUNT")
@AttributeOverride(name = "id", column = @Column(name = "DISCOUNT_ID"))
public class Discount extends AbstractAudit {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @OneToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "PRODUCT_ID")
    private Product product;

    @Column(name = "PERCENTAGE")
    private int percentage;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public Product getProduct() {
        return product;
    }

    public void setProduct(Product product) {
        this.product = product;
    }

    public int getPercentage() {
        return percentage;
    }

    public void setPercentage(int percentage) {
        this.percentage = percentage;
    }
}
