// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package catalog

//@todo routes to implement GH Issue#1
//POST   /V1/products
//PUT    /V1/products/:sku
//DELETE /V1/products/:sku
//GET    /V1/products
//GET    /V1/products/:sku
//GET    /V1/products/attributes/types
//GET    /V1/products/attributes/:attributeCode
//GET    /V1/products/attributes
//GET    /V1/categories/attributes/:attributeCode
//GET    /V1/categories/attributes
//GET    /V1/categories/attributes/:attributeCode/options
//POST   /V1/products/attributes
//PUT    /V1/products/attributes/:attributeCode
//DELETE /V1/products/attributes/:attributeCode
//GET    /V1/products/types
//GET    /V1/products/attribute-sets/sets/list
//GET    /V1/products/attribute-sets/:attributeSetId
//DELETE /V1/products/attribute-sets/:attributeSetId
//POST   /V1/products/attribute-sets
//PUT    /V1/products/attribute-sets/:attributeSetId
//GET    /V1/products/attribute-sets/:attributeSetId/attributes
//POST   /V1/products/attribute-sets/attributes
//DELETE /V1/products/attribute-sets/:attributeSetId/attributes/:attributeCode
//GET    /V1/products/attribute-sets/groups/list
//POST   /V1/products/attribute-sets/groups
//PUT    /V1/products/attribute-sets/:attributeSetId/groups
//DELETE /V1/products/attribute-sets/groups/:groupId
//GET    /V1/products/attributes/:attributeCode/options
//POST   /V1/products/attributes/:attributeCode/options
//DELETE /V1/products/attributes/:attributeCode/options/:optionId
//GET    /V1/products/media/types/:attributeSetName
//GET    /V1/products/:sku/media/:imageId
//POST   /V1/products/:sku/media
//PUT    /V1/products/:sku/media/:entryId
//DELETE /V1/products/:sku/media/:entryId
//GET    /V1/products/:sku/media
//GET    /V1/products/:sku/group-prices/
//POST   /V1/products/:sku/group-prices/:customerGroupId/price/:price
//DELETE /V1/products/:sku/group-prices/:customerGroupId/
//GET    /V1/products/:sku/group-prices/:customerGroupId/tiers
//POST   /V1/products/:sku/group-prices/:customerGroupId/tiers/:qty/price/:price
//DELETE /V1/products/:sku/group-prices/:customerGroupId/tiers/:qty
//DELETE /V1/categories/:categoryId
//GET    /V1/categories/:categoryId
//POST   /V1/categories
//GET    /V1/categories
//PUT    /V1/categories/:id
//PUT    /V1/categories/:categoryId/move
//GET    /V1/products/options/types
//GET    /V1/products/:sku/options
//GET    /V1/products/:sku/options/:optionId
//POST   /V1/products/options
//PUT    /V1/products/options/:optionId
//DELETE /V1/products/:sku/options/:optionId
//GET    /V1/products/links/types
//GET    /V1/products/links/:type/attributes
//GET    /V1/products/:sku/links/:type
//POST   /V1/products/:sku/links/:type
//DELETE /V1/products/:sku/links/:type/:linkedProductSku
//PUT    /V1/products/:sku/links/:link_type
//GET    /V1/categories/:categoryId/products
//POST   /V1/categories/:categoryId/products
//PUT    /V1/categories/:categoryId/products
//DELETE /V1/categories/:categoryId/products/:sku
