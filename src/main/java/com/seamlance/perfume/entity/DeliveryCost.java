package com.seamlance.perfume.entity;

import com.seamlance.perfume.audit.AbstractAudit;
import jakarta.persistence.*;

import java.math.BigDecimal;

@Entity
@Table(name = "DELIVERY_COST")
@AttributeOverride(name = "id", column = @Column(name = "DELIVERY_COST_ID"))
public class DeliveryCost extends AbstractAudit {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "LOCATION")
    private String location;

    @Column(name = "SERVICE")
    private String service;

    @Column(name = "COST")
    private BigDecimal cost;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getLocation() {
        return location;
    }

    public void setLocation(String location) {
        this.location = location;
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public BigDecimal getCost() {
        return cost;
    }

    public void setCost(BigDecimal cost) {
        this.cost = cost;
    }
}
