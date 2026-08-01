#!/bin/sh

set -eu


SERVICES_DIR="./pkg/services"
RESOURCE=$1
RESOURCE_PATH=$SERVICES_DIR/$RESOURCE
echo "Creating resource: $RESOURCE"

echo "PATH: $RESOURCE_PATH"
mkdir -v $RESOURCE_PATH
touch $RESOURCE_PATH/{routes,handlers,$RESOURCE,persistence,ui}.go

for file in routes handlers "$RESOURCE" persistence ui; do
	touch $RESOURCE_PATH/$file.go
  echo "package $RESOURCE" > "$RESOURCE_PATH/$file.go"
done

content=$(cat <<EOF
package $RESOURCE

import (
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

func SetupRoutes(r chi.Router, db bob.DB, ns *nats.Conn) error {
	svc := NewService(db, ns)
	handler := NewHandler(svc)

	r.Route("/$RESOURCE", func(r chi.Router) {
		r.Get("/", handler.${RESOURCE}Page)
	})

	return nil
}
EOF
)
echo "$content" > $RESOURCE_PATH/routes.go

content=$(cat <<EOF
package $RESOURCE

import (
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) ${RESOURCE}Page(w http.ResponseWriter, r *http.Request) {
}
EOF
)
echo "$content" > $RESOURCE_PATH/handlers.go

content=$(cat <<EOF
package $RESOURCE

import (
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db   bob.DB
	nats *nats.Conn
}

func NewService(db bob.DB, ns *nats.Conn) *Service {
	return &Service{db: db, nats: ns}
}
EOF
)
echo "$content" > $RESOURCE_PATH/${RESOURCE}.go
