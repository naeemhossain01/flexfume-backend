package com.seamlance.perfume.entity;

import com.fasterxml.jackson.annotation.JsonBackReference;
import com.fasterxml.jackson.annotation.JsonManagedReference;
import com.seamlance.perfume.audit.AbstractAudit;
import com.seamlance.perfume.constants.EntityConstant;
import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.*;

import java.util.List;

@Entity
@AllArgsConstructor
@RequiredArgsConstructor
@Data
@Table(name = "USER")
@AttributeOverride(name = "id", column = @Column(name = "USER_ID"))
public class User extends AbstractAudit {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "NAME")
    @NotBlank(message = EntityConstant.NAME_REQUIRED)
    private String name;

    @Column(name = "EMAIL", unique = true)
    private String email;

    @Column(name = "PHONE_NUMBER", unique = true)
    @NotBlank(message = EntityConstant.PHONE_NUMBER_REQUIRED)
    private String phoneNumber;

    @Column(name = "PASSWORD")
    @Size(min = 8)
    @NotBlank(message = EntityConstant.PASSWORD_REQUIRED)
    private String password;

    @Column(name = "ROLE", columnDefinition = "varchar(50) default USER")
    private String role;

    @OneToOne(mappedBy = "user", fetch = FetchType.LAZY, cascade = CascadeType.ALL)
    @JsonManagedReference
    private Address address;

    @OneToMany(mappedBy = "user", fetch = FetchType.LAZY, cascade = CascadeType.ALL)
    private List<Order> orders;

    @OneToMany(mappedBy = "user", fetch = FetchType.LAZY, cascade = CascadeType.ALL)
    private List<Cart> carts;

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

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getPhoneNumber() {
        return phoneNumber;
    }

    public void setPhoneNumber(String phoneNumber) {
        this.phoneNumber = phoneNumber;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getRole() {
        return role;
    }

    public void setRole(String role) {
        this.role = role;
    }

    public Address getAddress() {
        return address;
    }

    public List<Order> getOrders() {
        return orders;
    }

    public void setOrders(List<Order> orders) {
        this.orders = orders;
    }

    public void setAddress(Address address) {
        this.address = address;
    }

    public List<Cart> getCarts() {
        return carts;
    }

    public void setCarts(List<Cart> carts) {
        this.carts = carts;
    }
}
