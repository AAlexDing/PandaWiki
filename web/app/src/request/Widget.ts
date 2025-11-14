/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

import httpRequest, { ContentType, RequestParams } from "./httpClient";
import {
  DomainChatRequest,
  DomainChatSearchReq,
  DomainChatSearchResp,
  DomainResponse,
  PostShareV1ChatWidgetParams,
} from "./types";

/**
 * @description ChatWidget
 *
 * @tags Widget
 * @name PostShareV1ChatWidget
 * @summary ChatWidget
 * @request POST:/share/v1/chat/widget
 * @response `200` `DomainResponse` OK
 */

export const postShareV1ChatWidget = (
  query: PostShareV1ChatWidgetParams,
  request: DomainChatRequest,
  params: RequestParams = {},
) =>
  httpRequest<DomainResponse>({
    path: `/share/v1/chat/widget`,
    method: "POST",
    query: query,
    body: request,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description WidgetSearch
 *
 * @tags Widget
 * @name PostShareV1WidgetSearch
 * @summary WidgetSearch
 * @request POST:/share/v1/widget/search
 * @response `200` `(DomainResponse & {
    data?: DomainChatSearchResp,

})` OK
 */

export const postShareV1WidgetSearch = (
  request: DomainChatSearchReq,
  params: RequestParams = {},
) =>
  httpRequest<
    DomainResponse & {
      data?: DomainChatSearchResp;
    }
  >({
    path: `/share/v1/widget/search`,
    method: "POST",
    body: request,
    type: ContentType.Json,
    format: "json",
    ...params,
  });
