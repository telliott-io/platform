resource "kubernetes_cluster_role" "nginx_ingress_clusterrole" {
    depends_on = [kubernetes_namespace.ingress_nginx]
  
    metadata {
        name = "nginx-ingress-clusterrole"

        labels = {
        "app.kubernetes.io/name" = "ingress-nginx"

        "app.kubernetes.io/part-of" = "ingress-nginx"
        }
    }

    rule {
        verbs = [
            "get",
            "list",
            "watch",
        ]
        api_groups = [""]
        resources = [
            "services",
            "endpoints",
        ]
    }

    rule {
        api_groups = [
            "",
        ]
        resources = [
            "secrets",
        ]
        verbs = [
            "get",
            "list",
            "watch",
        ]
    }

    rule {
        api_groups = [
            "",
        ]
        resources = [
            "configmaps",
        ]
        verbs = [
            "get",
            "list",
            "watch",
            "update",
            "create",
        ]
    }

    rule {
        api_groups = [
            "",
        ]
        resources = [
            "pods",
        ]
        verbs = [
            "list",
            "watch",
        ]
    }

    rule {
        api_groups = [
            "",
        ]
        resources = [
            "namespaces",
        ]
        verbs = [
            "get",
            "list",
            "watch",
        ]
    }
    
    rule {
        api_groups = [
            "",
        ]
        resources = [
            "events",
        ]
        verbs = [
            "create",
            "patch",
            "list",
        ]
    }

    rule {
        api_groups = [
            "coordination.k8s.io",
        ]
        resources = [
            "leases",
        ]
        verbs = [
            "get",
            "list",
            "watch",
            "update",
            "create",
        ]
    }
    
    rule {
        api_groups = [
            "networking.k8s.io",
        ]
        resources = [
            "ingresses",
        ]
        verbs = [
            "list",
            "watch",
            "get",
        ]
    }
    
    rule {
        api_groups = [
            "networking.k8s.io",
        ]
        resources = [
            "ingresses/status",
        ]
        verbs = [
            "update",
        ]
    }

    rule {
        api_groups = [
            "k8s.nginx.org",
        ]
        resources = [
            "virtualservers",
            "virtualserverroutes",
            "globalconfigurations",
            "transportservers",
            "policies",
        ]
        verbs = [
            "list",
            "watch",
            "get",
        ]
    }
    
    rule {
        api_groups = [
            "k8s.nginx.org",
        ]
        resources = [
            "virtualservers/status",
            "virtualserverroutes/status",
            "policies/status",
            "transportservers/status",
            "dnsendpoints/status",
        ]
        verbs = [
            "update",
        ]
    }

    rule {
        api_groups = [
            "networking.k8s.io",
        ]
        resources = [
            "ingressclasses",
        ]
        verbs = [
            "get",
        ]
    }
    
    rule {
        api_groups = [
            "cis.f5.com",
        ]
        resources = [
            "ingresslinks",
        ]
        verbs = [
            "list",
            "watch",
            "get",
        ]
        }

    rule {
        api_groups = [
            "cert-manager.io",
        ]
        resources = [
            "certificates",
        ]
        verbs = [
            "list",
            "watch",
            "get",
            "update",
            "create",
            "delete",
        ]
    }
    
    rule {
        api_groups = [
            "externaldns.nginx.org",
        ]
        resources = [
            "dnsendpoints",
        ]
        verbs = [
            "list",
            "watch",
            "get",
            "update",
            "create",
            "delete",
        ]
    }
    
    rule {
        api_groups = [
            "externaldns.nginx.org",
        ]
        resources = [
            "dnsendpoints/status",
        ]
        verbs = [
            "update",
        ]
    }
}