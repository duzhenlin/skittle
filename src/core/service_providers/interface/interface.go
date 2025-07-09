// Package service_providers
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 10:49

package _interface

import "go.uber.org/dig"

type ServiceProviderInterface interface {
	Register(container *dig.Container) error
}
